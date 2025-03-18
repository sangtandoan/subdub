package chrono

import (
	"context"
	"fmt"
	"sync"
	"time"

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

func (c *chrono) checkSubscriptionsDailyToSendEmail() {
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

	for _, sub := range subs {
		jobs <- sub
	}

	cnt := 0
	for cnt != len(subs) {
		<-done
		cnt += 1
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

		c.checkSubscriptionsDailyToSendEmail()
	}
}
