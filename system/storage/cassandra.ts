import cassandra from "cassandra-driver";

const CASSANDRA_CLIENT_CONFIG = {
  contactPoints: process.env.OPENFOODDIARY_CASSANDRA_CONTACT_POINTS
    ? process.env.OPENFOODDIARY_CASSANDRA_CONTACT_POINTS.split(";")
    : ["localhost:9143"],
  localDataCenter:
    process.env.OPENFOODDIARY_CASSANDRA_LOCALDATACENTER ?? "datacenter1",
  credentials: {
    username: process.env.OPENFOODDIARY_CASSANDRA_USER ?? "cassandra",
    password: process.env.OPENFOODDIARY_CASSANDRA_PASSWORD ?? "cassandra",
  },
};

export const CASSANDRA_CLIENT = new cassandra.Client(CASSANDRA_CLIENT_CONFIG);
