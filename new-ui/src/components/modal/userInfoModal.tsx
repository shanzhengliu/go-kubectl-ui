import { Button, Label, Modal, TextInput } from "flowbite-react";
import { useUserStore } from "../../react-context/userNameContext";

export const UserInfoModal = (props: {
  show: boolean;
  setShow: (show: boolean) => void;
}) => {
  const { user } = useUserStore();
 


  return (
    <Modal show={props.show}  onClose={() => {
        props.setShow(false);
      }}>
      <Modal.Header>User Info</Modal.Header>

      <Modal.Body>
        <div className="max-w-md">
          <div className="mb-2 block">
            <Label htmlFor="username" value="User Name:" />
           <TextInput id="username" value={user?.name} readOnly={true} />
          </div>
          <div className="mb-2 block">
            <Label htmlFor="email" value="Email:" />
            <TextInput id="email" value={user?.email} readOnly={true} />
          </div>
           <div className="mb-2 block">
                <Label htmlFor="groups" value="Groups:" />
                <TextInput id="groups" value={user?.groups?.join(",")} readOnly={true} />
           </div>
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={() => props.setShow(false)}>Close</Button>
      </Modal.Footer>
    </Modal>
  );
};
