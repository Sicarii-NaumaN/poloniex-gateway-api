package prepare

import (
	"context"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"

	crontab "github.com/robfig/cron/v3"
)

const (
	everyMinute         = "* * * * *"
	everyFifteenMinutes = "*/15 * * * *"
)

func addCronHandlerFunc(
	ctx context.Context,
	c *crontab.Cron,
	schedule string,
	f func(ctx context.Context) error,
) {
	if schedule == "" {
		return
	}

	_, err := c.AddFunc(schedule, func() {
		inErr := f(ctx)
		if inErr != nil {
			logger.Errorf("cron handler func: %v", inErr)
		}
	})
	if err != nil {
		logger.Fatalf("schedule cron err: %v", err)
	}
}

func RunCronJobs(
	ctx context.Context,
) {

	c := crontab.New()
	defer c.Start()

	addCronHandlerFunc(ctx, c, everyMinute, nil)
}
