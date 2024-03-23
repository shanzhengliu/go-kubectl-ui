import { LOCALSHELL } from "../../utils/endpoints";


export const  DockerShell = ( props: {
    iframeKey: number;
} ) =>{
  return (
    <div  key={props.iframeKey} className="h-screen m-4">
     <iframe  className="w-full" src={LOCALSHELL} title="Docker Shell"  style={{ height: "calc(80vh - 2rem)" }} />
    </div>
  );
}