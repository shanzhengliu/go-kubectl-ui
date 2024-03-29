import { Button, FileInput, Label } from "flowbite-react"
import { useState } from "react";
import { OPENAPI_HELPER_UPLOAD } from "../../utils/endpoints";


export const  OpenAPIOnline = ()=>{
    const [file, setFile] = useState(null);

    const handleFileChange = (event:any) => {
        setFile(event.target.files[0]);
      };

    const uploadZip = async (file: any) => {
        const formData = new FormData();
        formData.append('file', file);
        try {
          const response = await fetch(OPENAPI_HELPER_UPLOAD, {
            method: 'POST',
            body: formData,
          });
          console.log(response);
        } catch (error) {
          console.error(error);
        }
      }  

    return (    
     <div className="flex">
        <div className="mb-2 block">
          <Label htmlFor="file-upload" value="Upload file" />
        </div>
        <FileInput accept=".zip"  className="flex-grow"  id="file-upload"  onChange={handleFileChange} />
        <Button className="ml-2" onClick={() => uploadZip(file)}>Submit</Button>
     </div>)
}