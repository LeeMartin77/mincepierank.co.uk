import cassandra from 'cassandra-driver';

const CASSANDRA_CLIENT_CONFIG = {
	contactPoints: process.env.CASSANDRA_CONTACT_POINTS
		? process.env.CASSANDRA_CONTACT_POINTS.split(';')
		: ['localhost:9143'],
	localDataCenter: process.env.CASSANDRA_LOCALDATACENTER ?? 'datacenter1',
	credentials: {
		username: process.env.CASSANDRA_USER ?? 'cassandra',
		password: process.env.CASSANDRA_PASSWORD ?? 'cassandra'
	}
};

export const CASSANDRA_CLIENT = new cassandra.Client(CASSANDRA_CLIENT_CONFIG);
