package notifier

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/itgcloud/gobackup/config"
	"github.com/itgcloud/gobackup/logger"
)

type Base struct {
	viper     *viper.Viper
	Name      string
	onSuccess bool
	onFailure bool
}

type Notifier interface {
	notify(title, message string) error
}

var (
	notifyTypeSuccess = 1
	notifyTypeFailure = 2
)

func newNotifier(name string, config config.SubConfig) (Notifier, *Base, error) {
	base := &Base{
		viper: config.Viper,
		Name:  name,
	}
	base.viper.SetDefault("on_success", true)
	base.viper.SetDefault("on_failure", true)

	base.onSuccess = base.viper.GetBool("on_success")
	base.onFailure = base.viper.GetBool("on_failure")

	switch config.Type {
	case "mail":
		mail, err := NewMail(base)
		return mail, base, err
	case "webhook":
		return NewWebhook(base), base, nil
	case "feishu":
		return NewFeishu(base), base, nil
	case "dingtalk":
		return NewDingtalk(base), base, nil
	case "discord":
		return NewDiscord(base), base, nil
	case "slack":
		return NewSlack(base), base, nil
	case "github":
		return NewGitHub(base), base, nil
	case "telegram":
		return NewTelegram(base), base, nil
	case "postmark":
		return NewPostmark(base), base, nil
	case "sendgrid":
		return NewSendGrid(base), base, nil
	case "ses":
		return NewSES(base), base, nil
	case "resend":
		return NewResend(base), base, nil
	}

	return nil, nil, fmt.Errorf("Notifier: %s is not supported", name)
}

func notify(model config.ModelConfig, title, message string, notifyType int) {
	logger := logger.Tag("Notifier")

	// remove common from notifiers
	newNotifiers := map[string]config.SubConfig{}
	for k, v := range model.Notifiers {
		newNotifiers[k] = v
	}
	delete(newNotifiers, "common")

	logger.Infof("Running %d Notifiers", len(model.Notifiers))
	for name, config := range newNotifiers {
		notifier, base, err := newNotifier(name, config)
		if err != nil {
			logger.Error(err)
			continue
		}

		if notifyType == notifyTypeSuccess {
			if base.onSuccess {
				if err := notifier.notify(title, message); err != nil {
					logger.Error(err)
				}
			}
		} else if notifyType == notifyTypeFailure {
			if base.onFailure {
				if err := notifier.notify(title, message); err != nil {
					logger.Error(err)
				}
			}
		}
	}
}

func Success(model config.ModelConfig) {
	title := fmt.Sprintf("[GoBackup] OK: Backup *%s* successful", model.Name)
	if model.Notifiers["common"].Viper.GetString("title_success") != "" {
		title = model.Notifiers["common"].Viper.GetString("title_success")
	}

	message := fmt.Sprintf("Backup of *%s* completed successfully at %s", model.Name, time.Now().Local())
	if model.Notifiers["common"].Viper.GetString("message_success") != "" {
		message = model.Notifiers["common"].Viper.GetString("message_success")
	}

	notify(model, title, message, notifyTypeSuccess)
}

func Failure(model config.ModelConfig, reason string) {
	title := fmt.Sprintf("[GoBackup] ERROR: Backup *%s* failed", model.Name)
	if model.Notifiers["common"].Viper.GetString("title_failure") != "" {
		title = model.Notifiers["common"].Viper.GetString("title_failure")
	}

	message := fmt.Sprintf("Backup of *%s* failed at %s:\n----------------------------------------------\n%s", model.Name, time.Now().Local(), reason)
	if model.Notifiers["common"].Viper.GetString("message_failure") != "" {
		message = model.Notifiers["common"].Viper.GetString("message_failure")
	}

	notify(model, title, message, notifyTypeFailure)
}
