
disable_snapshots()
allow_k8s_contexts(['test', 'ci'])

load('ext://min_tilt_version', 'min_tilt_version')
min_tilt_version('0.15.0') # includes fix for auto_init+False with tilt ci

include('./frontend/Tiltfile')
include('./backend/Tiltfile')
include('./notify/Tiltfile')
include('./shared/Tiltfile')
