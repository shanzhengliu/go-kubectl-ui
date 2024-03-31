import { LOCALSHELL } from "../../utils/endpoints";


export const  DockerShell = ( props: {
    iframeKey: number;
} ) =>{
  return (
    <div  key={props.iframeKey} className="h-screen m-4">
     <h2>this shell is the web shell for connecting to the docker container running this app, and it is a linux container which container "kubctl" and "helm" command, it means you can use cmd to talk to k8s cluster directly via this shell.</h2>
     <iframe  className="w-full h-full" src={LOCALSHELL} title="Docker Shell"  />
    </div>
  );
}