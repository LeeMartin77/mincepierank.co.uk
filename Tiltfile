load('ext://namespace', 'namespace_create', 'namespace_inject')

development_namespace='mincepierank-development'

namespace_create(development_namespace)

# Create external resources
k8s_yaml(namespace_inject(read_file('.development/cassandra.yml'), development_namespace))

# Create deployment
deployment = namespace_inject(read_file('.development/deployment.yml'), development_namespace)
k8s_yaml(deployment)

# Build: tell Tilt what images to build from which directories

docker_build('ghcr.io/leemartin77/mincepierank.co.uk', '.')

# Watch: tell Tilt how to connect locally (optional)

k8s_resource('cassandra', port_forwards="9143:9042", labels=["services"])
k8s_resource('mincepierank', port_forwards="4024:3000", labels=["application"],  resource_deps=['cassandra'])

local_resource('mincepierank local',
  serve_cmd='npm run dev',
  links=["http://localhost:3000"],
  serve_env={
    'CASSANDRA_CONTACT_POINTS': "localhost:9143",
    'CASSANDRA_USER':'cassandra',
    'CASSANDRA_PASSWORD':'cassandra',
  },
  dir='.',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["application"]
)

local_resource('migrate cassandra',
  cmd='npm run dev:migrations',
  dir='.',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tools"]
)
local_resource('seed cassandra',
  cmd='npm run dev:seedlocaltestdata',
  dir='.',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tools"]
)
