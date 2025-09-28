package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rabbitCmd = &cobra.Command{
	Use:   "rabbit",
	Short: "RabbitMQ management commands",
	Long:  `Manage RabbitMQ queues, exchanges, and bindings`,
}

var rabbitListQueuesCmd = &cobra.Command{
	Use:   "queues",
	Short: "List all queues",
	Run: func(cmd *cobra.Command, args []string) {
		listQueues()
	},
}

var rabbitCreateQueueCmd = &cobra.Command{
	Use:   "create-queue [queue-name]",
	Short: "Create a new queue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		durable, _ := cmd.Flags().GetBool("durable")
		autoDelete, _ := cmd.Flags().GetBool("auto-delete")
		createQueue(args[0], durable, autoDelete)
	},
}

var rabbitDeleteQueueCmd = &cobra.Command{
	Use:   "delete-queue [queue-name]",
	Short: "Delete a queue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteQueue(args[0])
	},
}

var rabbitPurgeQueueCmd = &cobra.Command{
	Use:   "purge [queue-name]",
	Short: "Purge all messages from a queue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		purgeQueue(args[0])
	},
}

var rabbitListExchangesCmd = &cobra.Command{
	Use:   "exchanges",
	Short: "List all exchanges",
	Run: func(cmd *cobra.Command, args []string) {
		listExchanges()
	},
}

var rabbitCreateExchangeCmd = &cobra.Command{
	Use:   "create-exchange [exchange-name]",
	Short: "Create a new exchange",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exchangeType, _ := cmd.Flags().GetString("type")
		durable, _ := cmd.Flags().GetBool("durable")
		autoDelete, _ := cmd.Flags().GetBool("auto-delete")
		createExchange(args[0], exchangeType, durable, autoDelete)
	},
}

var rabbitPublishCmd = &cobra.Command{
	Use:   "publish [exchange] [routing-key] [message]",
	Short: "Publish a message to an exchange",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		publishMessage(args[0], args[1], args[2])
	},
}

var rabbitStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show RabbitMQ statistics",
	Run: func(cmd *cobra.Command, args []string) {
		showStats()
	},
}

func init() {
	rabbitCmd.PersistentFlags().StringP("host", "H", "", "RabbitMQ host (env: RABBITMQ_HOST)")
	rabbitCmd.PersistentFlags().StringP("port", "p", "", "RabbitMQ management port (env: RABBITMQ_MANAGEMENT_PORT)")
	rabbitCmd.PersistentFlags().StringP("user", "u", "", "RabbitMQ user (env: RABBITMQ_DEFAULT_USER)")
	rabbitCmd.PersistentFlags().StringP("password", "P", "", "RabbitMQ password (env: RABBITMQ_DEFAULT_PASS)")
	rabbitCmd.PersistentFlags().StringP("vhost", "v", "", "RabbitMQ virtual host (env: RABBITMQ_DEFAULT_VHOST)")

	rabbitCreateQueueCmd.Flags().BoolP("durable", "d", true, "Make queue durable")
	rabbitCreateQueueCmd.Flags().BoolP("auto-delete", "a", false, "Auto-delete queue when unused")

	rabbitCreateExchangeCmd.Flags().StringP("type", "t", "direct", "Exchange type (direct, fanout, topic, headers)")
	rabbitCreateExchangeCmd.Flags().BoolP("durable", "d", true, "Make exchange durable")
	rabbitCreateExchangeCmd.Flags().BoolP("auto-delete", "a", false, "Auto-delete exchange when unused")

	rabbitCmd.AddCommand(rabbitListQueuesCmd)
	rabbitCmd.AddCommand(rabbitCreateQueueCmd)
	rabbitCmd.AddCommand(rabbitDeleteQueueCmd)
	rabbitCmd.AddCommand(rabbitPurgeQueueCmd)
	rabbitCmd.AddCommand(rabbitListExchangesCmd)
	rabbitCmd.AddCommand(rabbitCreateExchangeCmd)
	rabbitCmd.AddCommand(rabbitPublishCmd)
	rabbitCmd.AddCommand(rabbitStatsCmd)
}

func getRabbitMQURL(path string) string {
	host, _ := rabbitCmd.Flags().GetString("host")
	port, _ := rabbitCmd.Flags().GetString("port")

	if host == "" {
		host = viper.GetString("rabbitmq.host")
	}
	if port == "" {
		port = viper.GetString("rabbitmq.port")
	}

	return fmt.Sprintf("http://%s:%s/api%s", host, port, path)
}

func getRabbitMQAuth() (string, string) {
	user, _ := rabbitCmd.Flags().GetString("user")
	password, _ := rabbitCmd.Flags().GetString("password")

	if user == "" {
		user = viper.GetString("rabbitmq.user")
	}
	if password == "" {
		password = viper.GetString("rabbitmq.password")
	}

	return user, password
}

func makeRabbitMQRequest(method, path string, body interface{}) (*http.Response, error) {
	client := &http.Client{}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, getRabbitMQURL(path), reqBody)
	if err != nil {
		return nil, err
	}

	user, pass := getRabbitMQAuth()
	req.SetBasicAuth(user, pass)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return client.Do(req)
}

func listQueues() {
	vhost, _ := rabbitCmd.Flags().GetString("vhost")
	if vhost == "" {
		vhost = viper.GetString("rabbitmq.vhost")
	}
	resp, err := makeRabbitMQRequest("GET", fmt.Sprintf("/queues/%s", vhost), nil)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Red("Error: HTTP %d", resp.StatusCode)
		return
	}

	var queues []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&queues); err != nil {
		color.Red("Error parsing response: %v", err)
		return
	}

	color.Green("Queues in vhost '%s':", vhost)
	for _, queue := range queues {
		name := queue["name"].(string)
		messages := int(queue["messages"].(float64))
		consumers := int(queue["consumers"].(float64))
		fmt.Printf("  - %s (messages: %d, consumers: %d)\n", name, messages, consumers)
	}
}

func createQueue(queueName string, durable, autoDelete bool) {
	vhost, _ := rabbitCmd.Flags().GetString("vhost")
	if vhost == "" {
		vhost = viper.GetString("rabbitmq.vhost")
	}
	body := map[string]interface{}{
		"durable":     durable,
		"auto_delete": autoDelete,
	}

	resp, err := makeRabbitMQRequest("PUT", fmt.Sprintf("/queues/%s/%s", vhost, queueName), body)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		color.Green("Queue '%s' created successfully", queueName)
	} else {
		color.Red("Error creating queue: HTTP %d", resp.StatusCode)
	}
}

func deleteQueue(queueName string) {
	vhost, _ := rabbitCmd.Flags().GetString("vhost")
	if vhost == "" {
		vhost = viper.GetString("rabbitmq.vhost")
	}
	resp, err := makeRabbitMQRequest("DELETE", fmt.Sprintf("/queues/%s/%s", vhost, queueName), nil)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		color.Green("Queue '%s' deleted successfully", queueName)
	} else {
		color.Red("Error deleting queue: HTTP %d", resp.StatusCode)
	}
}

func purgeQueue(queueName string) {
	vhost, _ := rabbitCmd.Flags().GetString("vhost")
	if vhost == "" {
		vhost = viper.GetString("rabbitmq.vhost")
	}
	resp, err := makeRabbitMQRequest("DELETE", fmt.Sprintf("/queues/%s/%s/contents", vhost, queueName), nil)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		color.Green("Queue '%s' purged successfully", queueName)
	} else {
		color.Red("Error purging queue: HTTP %d", resp.StatusCode)
	}
}

func listExchanges() {
	vhost, _ := rabbitCmd.Flags().GetString("vhost")
	resp, err := makeRabbitMQRequest("GET", fmt.Sprintf("/exchanges/%s", vhost), nil)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Red("Error: HTTP %d", resp.StatusCode)
		return
	}

	var exchanges []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&exchanges); err != nil {
		color.Red("Error parsing response: %v", err)
		return
	}

	color.Green("Exchanges in vhost '%s':", vhost)
	for _, exchange := range exchanges {
		name := exchange["name"].(string)
		if name == "" {
			name = "(default)"
		}
		exchangeType := exchange["type"].(string)
		fmt.Printf("  - %s (type: %s)\n", name, exchangeType)
	}
}

func createExchange(exchangeName, exchangeType string, durable, autoDelete bool) {
	vhost, _ := rabbitCmd.Flags().GetString("vhost")
	body := map[string]interface{}{
		"type":        exchangeType,
		"durable":     durable,
		"auto_delete": autoDelete,
	}

	resp, err := makeRabbitMQRequest("PUT", fmt.Sprintf("/exchanges/%s/%s", vhost, exchangeName), body)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		color.Green("Exchange '%s' created successfully", exchangeName)
	} else {
		color.Red("Error creating exchange: HTTP %d", resp.StatusCode)
	}
}

func publishMessage(exchange, routingKey, message string) {
	vhost, _ := rabbitCmd.Flags().GetString("vhost")
	body := map[string]interface{}{
		"properties":       map[string]interface{}{},
		"routing_key":      routingKey,
		"payload":          message,
		"payload_encoding": "string",
	}

	resp, err := makeRabbitMQRequest("POST", fmt.Sprintf("/exchanges/%s/%s/publish", vhost, exchange), body)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		color.Green("Message published successfully")
	} else {
		color.Red("Error publishing message: HTTP %d", resp.StatusCode)
	}
}

func showStats() {
	resp, err := makeRabbitMQRequest("GET", "/overview", nil)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Red("Error: HTTP %d", resp.StatusCode)
		return
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		color.Red("Error parsing response: %v", err)
		return
	}

	color.Green("RabbitMQ Statistics:")

	if queueTotals, ok := stats["queue_totals"].(map[string]interface{}); ok {
		if messages, ok := queueTotals["messages"].(float64); ok {
			fmt.Printf("  Total Messages: %d\n", int(messages))
		}
		if messagesReady, ok := queueTotals["messages_ready"].(float64); ok {
			fmt.Printf("  Messages Ready: %d\n", int(messagesReady))
		}
		if messagesUnacknowledged, ok := queueTotals["messages_unacknowledged"].(float64); ok {
			fmt.Printf("  Messages Unacknowledged: %d\n", int(messagesUnacknowledged))
		}
	}

	if objectTotals, ok := stats["object_totals"].(map[string]interface{}); ok {
		if connections, ok := objectTotals["connections"].(float64); ok {
			fmt.Printf("  Connections: %d\n", int(connections))
		}
		if channels, ok := objectTotals["channels"].(float64); ok {
			fmt.Printf("  Channels: %d\n", int(channels))
		}
		if exchanges, ok := objectTotals["exchanges"].(float64); ok {
			fmt.Printf("  Exchanges: %d\n", int(exchanges))
		}
		if queues, ok := objectTotals["queues"].(float64); ok {
			fmt.Printf("  Queues: %d\n", int(queues))
		}
		if consumers, ok := objectTotals["consumers"].(float64); ok {
			fmt.Printf("  Consumers: %d\n", int(consumers))
		}
	}
}

func GetRabbitMQCommand() *cobra.Command {
	return rabbitCmd
}
