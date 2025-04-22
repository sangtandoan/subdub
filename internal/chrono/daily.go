package chrono

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sangtandoan/subscription_tracker/internal/pkg/enums"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/mailer"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
)

var days = [4]int{7, 5, 3, 1}

type Chrono interface {
	CheckSubscriptionDailyToSendEmail()
}

type chrono struct {
	subscriptionRepo repo.SubscriptionRepo
	userRepo         repo.UserRepo
	mailer           mailer.Mailer
}

func NewChrono(repo *repo.Repo, mailer mailer.Mailer) *chrono {
	return &chrono{subscriptionRepo: repo.Subscription, userRepo: repo.User, mailer: mailer}
}

func (c *chrono) CheckSubscriptionsDailyToSendEmail() {
	errsCh := make(chan error, len(days))

	wg := &sync.WaitGroup{}

	for _, num := range days {
		ctx := context.Background()
		wg.Add(1)
		go c.querySubsAtSpecifyNumDays(ctx, wg, num, errsCh)
	}

	wg.Wait()
	fmt.Println("wg wait done")
}

func (c *chrono) CheckSubscriptionsDailyToUpdateStartDate() {
	wg := &sync.WaitGroup{}

	subs, err := c.subscriptionRepo.GetSubscriptionsNeedUpdateStartAndEndDate(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan *repo.SubscriptionRow, 10)
	wg.Add(len(subs))

	done := func() {
		wg.Done()
	}

	for range 5 {
		go c.generateUpdateSubscriptionWorker(context.Background(), jobs, done)
	}

	for _, sub := range subs {
		jobs <- sub
	}

	wg.Wait()

	close(jobs)
}

func (c *chrono) generateUpdateSubscriptionWorker(
	ctx context.Context,
	jobs <-chan *repo.SubscriptionRow,
	done func(),
) {
	for job := range jobs {
		job.StartDate = job.EndDate
		duration, err := enums.ParseString2Duration(job.Duration)
		if err != nil {
			log.Fatal(err)
		}
		job.EndDate = duration.AddDurationToTime(job.StartDate)

		arg := repo.UpdateSubscriptionStartAndEndDateParams{
			ID:        job.ID,
			StartDate: job.StartDate,
			EndDate:   job.EndDate,
		}
		c.subscriptionRepo.UpdateSubscriptionStartAndEndDate(ctx, &arg)

		fmt.Println("Updated subscription with ID:", job.ID)
		done()
	}
}

func (c *chrono) querySubsAtSpecifyNumDays(
	ctx context.Context,
	wg *sync.WaitGroup,
	num int,
	errsCh chan<- error,
) {
	subs, err := c.subscriptionRepo.GetSubscriptionsBeforeNumDays(ctx, num)
	if err != nil {
		errsCh <- err
	}

	jobs := make(chan *repo.SubscriptionRow, 10)
	done := make(chan int, 10)

	c.generateWorkersPool(ctx, 3, num, jobs, done, errsCh)

	// instead we can use another goroutine to check for cnt,
	// and this will not block the main goroutine
	// and also will not be blocked by jobs channel
	// => deadlock will not occur
	cnt := 0
	go func() {
		for cnt != len(subs) {
			<-done
			cnt += 1
		}
	}()

	for _, sub := range subs {
		jobs <- sub
	}

	// if implement like this, deadlock can occur
	// because if all workers are busy not consunming jobs,
	// and subs is still more, then jobs will be blocked
	// and this code block will nerver excute
	// and if done channel full, worker can not send to done channel
	// so deadlock will occur
	//
	// cnt := 0
	// for cnt != len(subs) {
	// 	<-done
	// 	cnt += 1
	// }
	for cnt != len(subs) {
		continue
	}

	close(jobs)

	wg.Done()
}

func (c *chrono) generateWorkersPool(
	ctx context.Context,
	wokers int,
	numDays int,
	jobs <-chan *repo.SubscriptionRow,
	done chan<- int,
	errsChn chan<- error,
) {
	for range wokers {
		go c.sendEmail(ctx, numDays, jobs, done, errsChn)
	}
}

func (c *chrono) sendEmail(
	ctx context.Context,
	numDays int,
	jobs <-chan *repo.SubscriptionRow,
	done chan<- int,
	errsCh chan<- error,
) {
	// if jobs chan close, for loop will exit
	for job := range jobs {
		userID := job.UserID
		user, err := c.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			errsCh <- err
		}

		fmt.Printf("sending email to %s\n", user.Email)
		sendEmailReq := mailer.SendRequest{
			To:       []string{user.Email},
			Template: mailer.RemindTemplate,
			Data: mailer.RemindData{
				Name:        job.Name,
				NumDays:     numDays,
				Email:       user.Email,
				RenewalDate: job.EndDate,
			},
		}

		err = c.mailer.SendWithRetry(&sendEmailReq, 3)
		if err != nil {
			errsCh <- err
		}

		done <- 1
	}
}

func (c *chrono) ScheduleDailyTask(targetHour, targetMinute int) {
	// schedule daily excute at specify hour and minute
	for {
		now := time.Now()
		targetTime := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			targetHour,
			targetMinute,
			0,
			0,
			now.Location(),
		)

		if targetTime.Before(now) {
			targetTime = targetTime.AddDate(0, 0, 1)
		}

		waitDuration := targetTime.Sub(now)
		fmt.Println(waitDuration)
		time.Sleep(waitDuration)

		// c.CheckSubscriptionsDailyToSendEmail()
		c.CheckSubscriptionsDailyToUpdateStartDate()
	}
}
