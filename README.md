# Wal-transport
This is my final qualifying work for the undergraduate ITIS 2019-2023
This project is a Go-based data processing application that performs the following tasks:

1. Reads data from Write-Ahead Log (WAL).
2. Filters the data based on specific criteria.
3. Sends the filtered data to a message broker.
4. Contains a module that reads data from the message broker.
5. Converts the received data into a query.
6. Writes the converted query to a database.

## Features

- Provides customizable filtering options to process specific data.
- Integrates with a message broker for seamless data transfer.
- Converts received data into query format for database insertion.

## Prerequisites

To run this project, ensure you have the following installed on your system:

- Go (version 1.20 or higher)
- [Apache Kafka](https://kafka.apache.org/documentation/)
- [Nats](https://docs.nats.io/)
- [Database](https://www.postgresql.org/) 

## Getting Started

Follow the instructions below to get a local copy of the project up and running.

1. Clone the repository:

   ```bash
   git clone https://github.com/CAMELNINGA/WAL-transport.git
   ```

1. Change into the project directory:

   ```bash
   cd your-repository
   ```

1. Install project dependencies:

   ```bash
   go mod download
   ```

1. Configure the project:

   - Modify the configuration file (`deployments/config.yml`) to provide the necessary settings for your environment, including the message broker and database connection details.

1. Up modules copy_deamon and save_deamon:
   
    Example env copy_deamon in deployments/copy_deamon.env
    Example env save_deamon in deployments/save_deamon.env

   ```bash
    docker compose up copy_deamon save_deamon --build
   ```


1. Build the project:

   ```bash
   go build 
   ```

1. Run the data processing component:

   ```bash
   ./wal-transport send deployments/config.yml
   ```
   

**Note:** Ensure that the message broker and database are running and accessible.

## Contributing

Contributions are welcome! If you wish to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b my-feature-branch`.
3. Make your changes and commit them: `git commit -am 'Add new feature'`.
4. Push to the branch: `git push origin my-feature-branch`.
5. Submit a pull request.

Please ensure that your code follows the project's coding style and conventions.

## License

This project is licensed information in the [LICENSE](LICENSE) file.

## Contact

For any inquiries or feedback, please contact [rfvbkm0220@gmail.com](mailto:rfvbkm0220@gmail.com).