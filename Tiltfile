load('ext://namespace', 'namespace_create', 'namespace_inject')

config.define_bool("run-cypress")
config.define_bool("local")
cfg = config.parse()
run_cypress = cfg.get("run-cypress", False)
local = cfg.get("local", False)

development_namespace='mincepierank-development'

namespace_create(development_namespace)

local_images="%s/.development/testdata/images" %(os.getcwd())

# Create external resources
k8s_yaml(namespace_inject(read_file('.development/cassandra.yml'), development_namespace))
imgprssr = read_yaml_stream('.development/development-images.yaml')
for o in imgprssr:
  o['metadata']['namespace'] = development_namespace
imgprssr[3]['spec']['hostPath']['path'] = local_images
k8s_yaml(encode_yaml_stream(imgprssr))

# Create deployment
deployment = namespace_inject(read_file('kustomize/deployment.yml'), development_namespace)
k8s_yaml(deployment)
service = namespace_inject(read_file('kustomize/service.yml'), development_namespace)
k8s_yaml(service)

# Build: tell Tilt what images to build from which directories

docker_build('ghcr.io/leemartin77/mincepierank.co.uk', '.')

# Watch: tell Tilt how to connect locally (optional)

k8s_resource('cassandra', port_forwards="9145:9042", labels=["services"])
k8s_resource('mincepierank-imgprssr', port_forwards="3013:3013", labels=["services"])

k8s_resource('mincepierank', port_forwards="4025:3000", labels=["application"],  resource_deps=['cassandra'],
  auto_init=False if local else True,
  trigger_mode=TRIGGER_MODE_MANUAL if run_cypress or local else TRIGGER_MODE_AUTO)

local_resource('mincepierank local',
  serve_cmd='npm run dev',
  links=["http://localhost:5173"],
  serve_env={
    'CASSANDRA_CONTACT_POINTS': "localhost:9145",
    'CASSANDRA_USER':'cassandra',
    'CASSANDRA_PASSWORD':'cassandra',
    'READONLY': 'false',
    'AUTH_STRING': 'super-secret-string',
    'IMGPRSSR_DIR': local_images,
    'IMGPRSSR_ROOT': 'http://localhost:3013',
    'AUTHENTICATION': 'development'
  },
  dir='.',
  auto_init=True if local else False,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["application"]
)

local_resource('seed cassandra',
  cmd='npm run dev:seedlocaltestdata',
  dir='.',
  env={
    'CASSANDRA_CONTACT_POINTS': "localhost:9145",
    'CASSANDRA_USER':'cassandra',
    'CASSANDRA_PASSWORD':'cassandra',
    'READONLY': 'false',
    'AUTH_STRING': 'super-secret-string',
    'AUTHENTICATION': 'development'
  },
  auto_init=run_cypress,
  resource_deps=[] if local else ['mincepierank'],
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tools"]
)

local_resource('cypress run',
  env={'CYPRESS_BASE_URL':'http://localhost:4025'},
  cmd='npm run cypress:run',
  dir='.',
  resource_deps=['mincepierank', 'seed cassandra'],
  auto_init=run_cypress,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tests"]
)

local_resource('cypress open',
  env={'CYPRESS_BASE_URL':'http://localhost:4025'},
  serve_cmd='npm run cypress:open',
  serve_dir='.',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL,
  labels=["tests"]
)