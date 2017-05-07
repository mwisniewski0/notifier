package main

import (
  "io/ioutil"
  "fmt"
  "encoding/json"
  "time"
  "notifier/subscriptions"
  "notifier/emailutils"
)

type manager struct {
  sender *emailutils.EmailSender
  subscriptions []*subscriptions.Subscription
}

func shortest(durations []time.Duration) time.Duration {
  lowest := durations[0]
  for _, duration := range(durations) {
    if duration < lowest {
      lowest = duration
    }
  }
  return lowest
}

func (man *manager) run() {
  // This function only works if there is at least one subscription
  if (len(man.subscriptions) == 0) {
    fmt.Println("No active subscriptions")
    return
  }

  var lastChecked []time.Time

  now := time.Now()
  for i := 0; i < len(man.subscriptions); i++ {
    lastChecked = append(lastChecked, now)
  }

  timesTillNextCheck := make([]time.Duration, len(lastChecked))
  for {
    now = time.Now()
    for index, subscr := range(man.subscriptions) {
      timesTillNextCheck[index] = subscr.CheckInterval() - (now.Sub(lastChecked[index]))
    }

    nextCheckStartsIn := shortest(timesTillNextCheck)
    if nextCheckStartsIn <= 0 {
      // There are overdue checks, process them
      for index, timeTillNextCheck := range(timesTillNextCheck) {
        if timeTillNextCheck <= 0 {
          // TODO: fix the below
          go man.subscriptions[index].Check(man.sender)
          lastChecked[index] = now
        }
      }
    } else {
      // wait for the next check
      time.Sleep(nextCheckStartsIn)
    }
  }
}

type ManagerConfig struct {
  Sender emailutils.SenderConfig
  Subscriptions []subscriptions.SubscriptionDescriptor
}

func (config *ManagerConfig) makeManager() (*manager, error) {
  newManager := new(manager)
  newManager.sender = config.Sender.MakeSender()

  for _, descriptor := range(config.Subscriptions) {
    newSubscription, err := descriptor.MakeSubscription()
    if err != nil {
      return nil, err
    }
    newManager.subscriptions = append(newManager.subscriptions, newSubscription)
  }

  return newManager, nil
}

func loadManager(path string) (*manager, error) {
  contents, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }

  config := new(ManagerConfig)
  err = json.Unmarshal(contents, config)

  if err != nil {
    return nil, err
  }

  newManager, err := config.makeManager()
  return newManager, err
}

func main() {
  manager, err := loadManager("config.json")
  if err != nil {
    fmt.Println("Could not parse the JSON configuration.")
    fmt.Println(err)
    return
  }

  manager.run()
}
