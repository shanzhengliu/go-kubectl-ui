import { Button,  Textarea } from "flowbite-react"
import { inputHook } from "../../hooks/inputhook";
import Swal from "sweetalert2";


export const JsonFormatter =()=>{
    const [json,setJson, onChangeJson] = inputHook("");
    const [formattedJson, setFormattedJson] = inputHook("");
    const formatJsonRun = ()=>{
        try {
            setFormattedJson(JSON.stringify(JSON.parse(json), null, 2));
        } catch (error) {
            setFormattedJson("Invalid Json");
        }
    }

    const copyToClipboard = ()=>{
        navigator.clipboard.writeText(formattedJson);
        Swal.fire({
            icon: 'success',
            title: 'Copied to clipboard',
            showConfirmButton: false,
            timer: 1500
          })

    }

    const compressJson = ()=>{
        try {
            setFormattedJson(JSON.stringify(JSON.parse(json)));
        } catch (error) {
            setFormattedJson("Invalid Json");
        }
    }

    const clearJson = ()=>{
        setFormattedJson("");
        setJson("")
    }

    
    return (
        <div  >
            <div className="flex h-screen">
            <Textarea className="w-full resize-none m-4" placeholder="Enter Json here" value={json} onChange={onChangeJson}  />
                    <div className="block mt-24">
                       <Button className="w-24 h-12 mt-4 " gradientDuoTone="purpleToBlue" onClick={formatJsonRun}>Format</Button>
                       <Button className="w-24 h-12 mt-4 " gradientDuoTone="purpleToPink" onClick={compressJson}>Compress</Button>
                       <Button className="w-24 h-12 mt-4 " gradientDuoTone="pinkToOrange" onClick={clearJson}>Clear</Button>
                       <Button className="w-24 h-12 mt-4 " color="success" onClick={copyToClipboard}>Copy</Button>
                    </div>
            <Textarea className="w-full resize-none m-4" placeholder="Formatted Json" value={formattedJson}   readOnly/>
           </div>
        </div>
    )
}