package subscriptions

import (
  "encoding/json"
  "errors"
  "sync"
  "time"
  "fmt"
  "notifier/emailutils"
  "strings"
)

type NotificationChecker interface {
  // NextCheck checks if a notification should be sent. If so, it returns the
  // contents of the email to be sent. Otherwise, it returns nil.
  NextCheck() (*string, error)

  // IsConcurrent indicates whether nextCheck() can work properly in a
  // concurrent setting. If there can only be one instance of nextCheck()
  // running, isConcurrent should return false. Otherwise, it should return true
  IsConcurrent() bool
}

type Subscription struct {
  // checkInterval indicates what should be the time between consecutive checks
  // performed by this subscription
  checkInterval time.Duration

  // checker represents the checking logic of the subscription. Every time
  // checkInterval elapses, checker.nextCheck() will be called to determine if
  // a notification is necessary.
  checker NotificationChecker

  // recipients represents the email addresses of all people subscribed to this
  // subscription.
  recipients []string

  wg sync.WaitGroup
}

func (sub *Subscription) CheckInterval() time.Duration {
  return sub.checkInterval
}

func (sub *Subscription) Checker() *NotificationChecker {
  return &sub.checker
}

func (sub *Subscription) Recipients() []string {
  return sub.recipients
}

type SubscriptionDescriptor struct {
  // CheckInterval is the amount of milliseconds between consecutive checks
  CheckInterval int64

  // CheckerType specifies the type of the checker to be read from JSON
  CheckerType string

  Recipients []string

  // CheckerData is the JSON data that will later be parsed into the appropriate
  // notificationChecker
  CheckerData json.RawMessage
}


func (sub *Subscription) Check(sender *emailutils.EmailSender) {
  if !sub.checker.IsConcurrent() {
    sub.wg.Wait()
  }
  sub.wg.Add(1)
  go sub.processCheck(sender)
}

func (sub *Subscription) processCheck(sender *emailutils.EmailSender) {
  defer sub.wg.Done()

  message, err := sub.checker.NextCheck()
  if err != nil {
    fmt.Println("Notificaiton check failed.")
    fmt.Println(err)
  } else if(message != nil) {
    for _, recipient := range(sub.recipients) {
      err := sender.Send(*message, recipient)
      if err != nil {
        fmt.Println("Problems sending the email")
        fmt.Println(err)
      }
    }
    fmt.Println("Sent notification emails to: " + strings.Join(sub.recipients, ", "))
  }
}

func (descriptor *SubscriptionDescriptor) MakeSubscription() (*Subscription, error) {
  newSubscription := new(Subscription)
  newSubscription.checkInterval = time.Duration(descriptor.CheckInterval) * time.Millisecond
  newSubscription.recipients = descriptor.Recipients

  // TODO: come up with a better design for this:
  switch descriptor.CheckerType {
  case "MyLeaseStarVacancy":
    newChecker := new(MyLeaseStarVacancy)
    err := json.Unmarshal(descriptor.CheckerData, newChecker)
    if err != nil {
      return nil, err
    }
    newSubscription.checker = newChecker

  default:
    return nil, errors.New("Unsupported notification type")
  }

  return newSubscription, nil
}
