function bulkedit(){
  var elements = document.getElementsByClassName('bulk-edit')
  for (var i = 0; i < elements.length; i++){
        elements[i].style.display = 'table-cell';
  }
  document.getElementById('bulkedit').style.display ='none';
  document.getElementById('activebulkedit').style.display ='block';
  document.getElementById('submitbulkedit').style.display ='block';
}

function nonbulkedit(){
  var elements = document.getElementsByClassName('bulk-edit')
  for (var i = 0; i < elements.length; i++){
        elements[i].style.display = 'none';
  }
  document.getElementById('bulkedit').style.display ='block';
  document.getElementById('activebulkedit').style.display ='none';
  document.getElementById('submitbulkedit').style.display ='none';
}

function hideended(){
  document.getElementById('hideended').style.display ='none';
  document.getElementById('showended').style.display ='block';
  var input, filter, table, tr, td, i;
  input = "✘";
  table = document.getElementById("showstable");
  tr = table.getElementsByTagName("tr");
  for (i = 0; i < tr.length; i++) {
    td = tr[i].getElementsByTagName("td")[2];
    if (td) {
      if (td.innerHTML.indexOf(input) > -1) {
        tr[i].style.display = "none";
      } else {
        tr[i].style.display = "";
      }
    }
  }
}

function showended(){
  document.getElementById('hideended').style.display ='block';
  document.getElementById('showended').style.display ='none';
  var input, filter, table, tr, td, i;
  input = "✘";
  table = document.getElementById("showstable");
  tr = table.getElementsByTagName("tr");
  for (i = 0; i < tr.length; i++) {
    td = tr[i].getElementsByTagName("td")[2];
    if (td) {
      if (td.innerHTML.indexOf(input) > -1) {
        tr[i].style.display = "";
      } else {
        tr[i].style.display = "";
      }
    }
  }
}

