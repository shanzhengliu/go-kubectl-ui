(function(){
    const password = localStorage.getItem("password")
    fetch("/api/login?password="+password).then((result)=>{
        if(result.status != 200){
            window.location.href = "/auth"
        }
    })
}())