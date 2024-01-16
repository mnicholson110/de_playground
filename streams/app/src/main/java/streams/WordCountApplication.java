package streams;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Produced;
import org.apache.kafka.streams.kstream.Printed;
import org.json.JSONObject;

import java.util.Properties;

public class WordCountApplication {
    public static void main(String[] args) {
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "order-filter-application");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9094");
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.Integer().getClass().getName());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass().getName());

        StreamsBuilder builder = new StreamsBuilder();
        KStream<Integer, String> orderStream = builder.stream("orderTopic", Consumed.with(Serdes.Integer(), Serdes.String()));

        KStream<Integer, String> filteredOrders = orderStream
                .filter((key, value) -> {
                    // Parse the JSON value and check the "OrderStatus" field
                    String orderStatus = parseOrderStatus(value);
                    return "delivered".equals(orderStatus);
                });

        filteredOrders.to("filteredOrderTopic", Produced.with(Serdes.Integer(), Serdes.String()));

        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();
    }

    private static String parseOrderStatus(String jsonValue) {
        try {
            // Parse the JSON value using org.json
            JSONObject jsonObject = new JSONObject(jsonValue);
            // Assuming the JSON structure is like: {"OrderStatus": "delivered", ...}
            return jsonObject.optString("OrderStatus", ""); // Get the "OrderStatus" field
        } catch (Exception e) {
            // Handle JSON parsing exception
            return "";
        }
    }
}

