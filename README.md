## **Log Ingestor**


**Log Ingestor**
A simple logging service that reports the error inside a http server. It uses `Kafka` for message queuing and `MongoDB` for storing the logs. 


### **Features**
* Simulates HTTP server errors.
* Publishes error logs to Kafka.
* Consumes error logs from Kafka.
* Stores error logs in MongoDB.
* Simple configuration for Kafka and MongoDB connection.

### **Prerequisites**
Before running the project, ensure you have the following installed:

* Go (version 1.18 or higher)
* Docker (for running Kafka and MongoDB in containers)

### **Usage**
##### 1. Clone the repository
```sh
https://github.com/PrantaDas/log-ingestor.git
```

##### 2. Change directory
```sh
cd log-ingestor
```

##### 3. Run docker compose file
```sh
docker compose up -d
```

##### 4. Verify kafka and zookeeper
```shell
nc -zv localhost 9092
nc -zv localhost 2181
```

##### 5. Push log to the message queue ( Send a request to the url to simulate error )
```sh
http://localhost:8080/simulate-error
```

### **Contributing**
Contributions are welcome! If you have suggestions for improvements or features, please create an issue or submit a pull request.

### **License**
This project is licensed under the MIT License. See the [LICENSE](https://github.com/PrantaDas/log-ingestor/blob/main/LICENSE) file for details.
