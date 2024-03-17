const backend = import.meta.env.VITE_BACKEND_URL
export const CONFIGMAP=backend+'/configmap'
export const DEPLOYMENT=backend+'/deployment'
export const INGRESS=backend+'/ingress'
export const POD=backend+'/pod'
export const SERVICE=backend+'/service'
export const PODLOGS=backend+'/api/podLogs'
export const WEBSHELL=backend+'/webshell'
export const PODYAML=backend+'/api/podYaml'
export const LOCALSHELL=backend+'/localshell'
export const CONFIGMAP_DETAIL=backend+'/api/configmap-detail'