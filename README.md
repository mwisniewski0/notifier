# notifier

A simple project I wrote in a couple of hours to learn Golang. Notifier is able to notify users via email about various events. Currently it has a built-in event for new vacancies in apartment complexes with websites hosted on MyLeaseStar (I used it to get myself a nice apartment in Colorado Springs :)).

The system is very flexible, and one can write code that can check for practically any event. The notifier will make sure to call your checking function in specified intervals of time. Then - if your checking function returns an email message - it will send it out to the users subscribed to the event.
