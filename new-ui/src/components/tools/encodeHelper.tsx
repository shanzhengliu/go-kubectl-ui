import { Button, Textarea } from "flowbite-react";
import { inputHook } from "../../hooks/inputhook";
import Swal from "sweetalert2";

export const EncodeHelper = () => {
  const [input, setInput, onChangeInput] = inputHook("");
  const [output, setOutput] = inputHook("");
  const copyToClipboard = () => {
    navigator.clipboard.writeText(output);
    Swal.fire({
      icon: "success",
      title: "Copied to clipboard",
      showConfirmButton: false,
      timer: 1500,
    });
  };

  const clearInput = () => {
    setInput("");
    setOutput("");
  };

  const base64Encode = () => {
    try {
      setOutput(btoa(input));
    } catch (error) {
      setOutput("Invalid Text");
    }
  };

  const based64Decode = () => {
    try {
      setOutput(atob(input));
    } catch (error) {
      setOutput("Invalid Text");
    }
  };

  const urlEncode = () => {
    try {
      setOutput(encodeURIComponent(input));
    } catch (error) {
      setOutput("Invalid Text");
    }
  };

  const urlDecode = () => {
    try {
      setOutput(decodeURIComponent(input));
    } catch (error) {
      setOutput("Invalid Text");
    }
  };

  return (
    <div>
      <div className="flex h-screen">
        <Textarea
          className="w-full resize-none m-4"
          placeholder="Input Text"
          value={input}
          onChange={onChangeInput}
        />
        <div className="block mt-24">
          <Button
            className="w-24 h-12 mt-4 "
            gradientDuoTone="purpleToBlue"
            onClick={base64Encode}
          >
            Base64 Encode
          </Button>
          <Button
            className="w-24 h-12 mt-4 "
            gradientDuoTone="purpleToPink"
            onClick={based64Decode}
          >
            Base64 Decode
          </Button>
          <Button
            className="w-24 h-12 mt-4 "
            gradientDuoTone="pinkToOrange"
            onClick={urlEncode}
          >
            Url Encode
          </Button>
          <Button
            className="w-24 h-12 mt-4 "
            gradientDuoTone="tealToLime"
            onClick={urlDecode}
          >
            Url Decode
          </Button>
          <Button
            className="w-24 h-12 mt-4 "
            gradientDuoTone="cyanToBlue"
            onClick={clearInput}
          >
            Clear
          </Button>
          <Button
            className="w-24 h-12 mt-4 "
            color="success"
            onClick={copyToClipboard}
          >
            Copy
          </Button>
        </div>
        <Textarea
          className="w-full resize-none m-4"
          placeholder="Output Text"
          value={output}
          readOnly
        />
      </div>
    </div>
  );
};
