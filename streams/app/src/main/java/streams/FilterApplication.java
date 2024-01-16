package streams;

import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Produced;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.json.JSONObject;
import java.util.Properties;

public class FilterApplication {
  public static void main(String[] args) throws Exception {
    Properties props = new Properties();
    props.put(StreamsConfig.APPLICATION_ID_CONFIG, "filter-application");
    props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9094");
    props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
    props.put(StreamsConfig.REPLICATION_FACTOR_CONFIG, -1);
//    props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.Integer().getClass().getName());
 //   props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass().getName());
    System.out.println("Starting filter application");
    StreamsBuilder builder = new StreamsBuilder();
    builder.stream("orderTopic", Consumed.with(Serdes.Integer(), Serdes.Bytes()))
      .peek((key, value) -> System.out.println("Key: " + key + " Value: " + value));
      //.filter((key, value) -> parseOrderStatus(value).equals("delivered"))
      //.to("filteredTopic", Produced.with(Serdes.Integer(), Serdes.Bytes()));

    KafkaStreams streams = new KafkaStreams(builder.build(), props);
    streams.start();
  }

  private static String parseOrderStatus(Object jsonValue) {
    try {
      System.out.println("Parsing JSON");
      JSONObject jsonObject = new JSONObject(jsonValue);
      return jsonObject.optString("OrderStatus", "");
    } catch (Exception e) {
      System.out.println("Error parsing JSON");
      return "";
    }
  }
}

