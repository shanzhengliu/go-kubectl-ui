{{define "navigator"}}
  <script>
        
        function search(value){
            var tr = document.querySelectorAll('.tbbody')
            tr.forEach(element => {
              var td =  element.querySelectorAll('td');
              const found = [...td].find(cell => {
                
                    return cell.textContent.includes(value);
              });
              if( !found){
                element.style.display = "none"
              }
              else{
                element.style.display = "table-row"
              }
            });
        }
      </script>
<a href="/pod">Pod</a>
<a href="/deployment">Deployment</a>
<a href="/configmap">Configmap</a>
<a href="/ingress">Ingress</a>
 <br/>
search <input  id=search type="text" onchange="search(value)"/>
{{end}}