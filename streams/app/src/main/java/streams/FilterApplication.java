package streams;

import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.KStream;
import org.json.JSONObject;
import java.util.Properties;

public class FilterApplication {
  public static void main(String[] args) throws Exception {
    Properties props = new Properties();
    props.put(StreamsConfig.APPLICATION_ID_CONFIG, "application");
    props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9094");
    props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.Integer().getClass().getName());
    props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass().getName());
    
    StreamsBuilder builder = new StreamsBuilder();
    builder.stream("orderTopic")
      .filter((key, value) -> parseOrderStatus(value).equals("delivered"))
      .to("filteredTopic");

    KafkaStreams streams = new KafkaStreams(builder.build(), props);

    streams.start();
  }

  private static String parseOrderStatus(Object jsonValue) {
    try {
      JSONObject jsonObject = new JSONObject(jsonValue);
      return jsonObject.optString("OrderStatus", "");
    } catch (Exception e) {
      return "";
    }
  }
}

