package main

import waltransport "github.com/CAMELNINGA/WAL-transport.git/cmd/wal-transport"

func main() {
	waltransport.Execute()
}

// // CLi chicken config and send to kafka
// func main() {
// 	ctx := context.Background()

// 	cli := config.Cli{}
// 	f, err := os.Open("diplom/config.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	fmt.Println("The File is opened successfully...")

// 	if err := cli.Parse(f); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("The file is parsed successfully...")

// 	if len(cli.Kafka.Brokers) == 0 {
// 		panic("kafka brokers not set")
// 	}

// 	var flag kafka.Bits
// 	flag = kafka.Set(flag, kafka.Producer)
// 	k := kafka.NewKafka(
// 		kafka.WithBrokers(cli.Kafka.Brokers),
// 		kafka.WithTopic(cli.Kafka.Topic),
// 		kafka.WithGroupID(cli.Kafka.GroupID),
// 		kafka.WithContext(ctx),
// 		kafka.WithFlags(flag),
// 	)

// 	for name, v := range cli.Deamons {

// 		k.PublishConfig(name, v)
// 	}

// }
