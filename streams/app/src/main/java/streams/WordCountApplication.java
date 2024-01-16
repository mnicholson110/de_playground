package streams;

import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.KTable;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.kstream.Produced;
import org.apache.kafka.streams.state.KeyValueStore;

import java.util.Arrays;
import java.util.Properties;

public class WordCountApplication {

   public static void main(final String[] args) throws Exception {
       Properties props = new Properties();
       props.put(StreamsConfig.APPLICATION_ID_CONFIG, "wordcount-application");
       props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9094");
       props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.Json().getClass());
       props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.Json().getClass());

       StreamsBuilder builder = new StreamsBuilder();
       KStream<Int, Int> records = builder.stream("orderTopic");
       KTable<Int, Int> maxStatus = records
               .groupBy((key, value) -> KeyValue.pair(value, value))
               .reduce((aggValue, newValue) -> newValue);
       maxStatus.toStream().to("orderMaxTopic", Produced.with(Serdes.Int(), Serdes.Int()));

       KafkaStreams streams = new KafkaStreams(builder.build(), props);
       streams.start();
   }

}
