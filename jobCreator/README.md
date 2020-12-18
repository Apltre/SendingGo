# SendingGo
1)jobCreator is a web API for pushing new "jobs" into MongoDB
2)sendingQueue is a service that takes new "jobs" from DB. Thеn deliver them to sendingService or sendingResult service through RabbitMq.
3)sendingService sends data to outer services and returns operation state back to sendingQueue through RabbitMq
4)sendingResult does some work like sending notifications to inner systems depending on "job" sending state