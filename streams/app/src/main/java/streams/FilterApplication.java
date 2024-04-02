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
    StreamsBuilder builder = new StreamsBuilder();
    builder.stream("orderTopic", Consumed.with(Serdes.Bytes(), Serdes.String()))
      .filter((key, value) -> parseOrderStatus(value).equals("delivered"))
      .to("deliveredTopic", Produced.with(Serdes.Bytes(), Serdes.String()));

    KafkaStreams streams = new KafkaStreams(builder.build(), props);
    streams.start();
    Runtime.getRuntime().addShutdownHook(new Thread(streams::close));
  }

  private static String parseOrderStatus(String jsonValue) {
    try {
      JSONObject jsonObject = new JSONObject(jsonValue);
      return jsonObject.optString("orderStatus", "");
    } catch (Exception e) {
      return "";
    }
  }
}

