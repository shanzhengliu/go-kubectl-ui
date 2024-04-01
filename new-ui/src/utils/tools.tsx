import Swal from "sweetalert2";

export const errorHander = (error: any) => {
    Swal.fire({
        icon: "error",
        text: error.response.data.error,
    });
}