load('ext://namespace', 'namespace_create', 'namespace_inject')

config.define_bool("run-cypress")
cfg = config.parse()
run_cypress = cfg.get("run-cypress", False)

development_namespace='mincepierank-development'

namespace_create(development_namespace)

# Create external resources
k8s_yaml(namespace_inject(read_file('.development/cassandra.yml'), development_namespace))

# Create deployment
deployment = namespace_inject(read_file('.development/deployment.yml'), development_namespace)
k8s_yaml(deployment)
deployment_svelte = namespace_inject(read_file('.development/deployment.svelte.yml'), development_namespace)
k8s_yaml(deployment_svelte)

# Build: tell Tilt what images to build from which directories

docker_build('ghcr.io/leemartin77/mincepierank.co.uk', '.')
docker_build('ghcr.io/leemartin77/mincepierank.co.uk-svelte', 'svelte')

# Watch: tell Tilt how to connect locally (optional)

k8s_resource('cassandra', port_forwards="9143:9042", labels=["services"])
k8s_resource('mincepierank', port_forwards="4024:3000", labels=["application"],  resource_deps=['cassandra'],
  auto_init=run_cypress,
  trigger_mode=TRIGGER_MODE_MANUAL if run_cypress else TRIGGER_MODE_AUTO)

k8s_resource('mincepierank-svelte', port_forwards="4025:3000", labels=["application"],  resource_deps=['cassandra'],
  auto_init=run_cypress,
  trigger_mode=TRIGGER_MODE_MANUAL if run_cypress else TRIGGER_MODE_AUTO)

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

local_resource('mincepierank svelte local',
  serve_dir='svelte',
  serve_cmd='npm run dev',
  links=["http://localhost:5173"],
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
  auto_init=run_cypress,
  resource_deps=['cassandra'],
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tools"]
)
local_resource('seed cassandra',
  cmd='npm run dev:seedlocaltestdata',
  dir='.',
  auto_init=run_cypress,
  resource_deps=['migrate cassandra'],
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tools"]
)

local_resource('cypress run',
  env={'CYPRESS_BASE_URL':'http://localhost:4024'},
  cmd='npm run cypress:run',
  dir='.',
  resource_deps=['mincepierank', 'seed cassandra'],
  auto_init=run_cypress,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tests"]
)

local_resource('cypress open - local',
  env={'CYPRESS_BASE_URL':'http://localhost:3000'},
  serve_cmd='npm run cypress:open',
  serve_dir='.',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tests"]
)