var InfrastructuresTable = {
  props: ['context'],
  data() {
    return {
      infrastructures: 'loading',
      running: 'loading',
      templates: 'loading',
      infrastructureToLoad: '',
      templateToLoad: '',
      infraToEdit: 'test',
      currentStatus: '',
      testOutput: ''
    }
  },
  template: `
    <div class="row">
      <div class="col-lg-2">
        <div class="table-responsive">
          <table class="table table-hover">
            <caption>Tests</caption>
            <tbody v-for="infrastructure in infrastructures" v-bind:key="infrastructure.Name">
              <tr>
                <td v-on:click="setInfrastructure(infrastructure)">{{ infrastructure.Name }}</td>
                <button href="#" data-bs-toggle="modal" :data-bs-target="'#modal-' + infrastructure.Name" v-on:click="updateInfraToEdit(infrastructure.Name)"> edit </button>
                <div class="modal fade" :id="'modal-' + infrastructure.Name" tabindex="-1" aria-hidden="true">
                  <div class="modal-dialog modal-dialog-scrollable modal-xl">
                    <div class="modal-content">
                      <div class="modal-header">
                        <h6 class="modal-title"><strong>Test Name:</strong><em> {{ infrastructure.Name }} </em> </h6>
                        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                      </div>
                      <div class="modal-body">
                        <form >
                          <div class="mb-3">
                            <textarea class="form-control" :bind="infraToEdit" :id="'message-text-'+ infrastructure.Name">{{ infraToEdit }}</textarea>
                          </div>
                        </form>
                      </div>
                      <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                        <button type="button" class="btn btn-danger" v-on:click="deleteInfra(infrastructure.Name)" data-bs-dismiss="modal">Delete</button>
                        <button type="button" class="btn btn-primary" v-on:click="sendInfraToEdit(infrastructure.Name)" data-bs-dismiss="modal">Save</button>
                      </div>
                    </div>
                  </div>
                </div>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="row" style="margin-bottom: 15px;">
          <div class="col-8">
            <input type="text" id="textNewInfra" class="form-control" placeholder="Test Name" aria-label="Test name" style="padding:0rem 0rem; font-size:inherit;">
          </div>
          <div class="col">
            <button v-on:click="sendNewInfra()">Add</button>
          </div>
        </div>

        <div class="table-responsive">
          <table class="table table-hover">
            <caption>Templates</caption>
            <tbody v-for="template in templates" v-bind:key="template.Name">
              <tr>
                <td v-on:click="setTemplate(template)">{{ template.Name }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="col" id="infrastructure-zone">
      </div>
     
      <div class="modal fade" :id="'modal-show-test-results'" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog modal-xl">
          <div class="modal-content">
            <div class="modal-header">
              <h6 class="modal-title"><strong>Test Output:</strong></h6>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body" style="white-space: pre-line"><code id="modal-show-test-results-content"></code></div>
          </div>
        </div>
      </div>

    </div>
  `,
  created() {
    this.searchInfrastructure();
    this.refresh();
    this.searchTemplates();
  },
  methods: {
    searchInfrastructure: function () {
      fetch('/api/infrastructures')
        .then(response => response.json())
        .then(json => {this.infrastructures = json.resultList.result;});
    },
    searchRunning: function () {
      fetch('/api/running')
        .then(response => response.json())
        .then(json => {this.running = json.resultList.result;});
    },
    searchTemplates: function () {
      fetch('/api/templates')
        .then(response => response.json())
        .then(json => {this.templates = json.resultList.result;});
    },
    setInfrastructure: function(infrastructure){
      this.templateToLoad = '';
      this.infrastructureToLoad = infrastructure;
      this.loadInfrasctructure(this.infrastructureToLoad);
    },
    setTemplate: function(template){
      this.infrastructureToLoad = '';
      this.loadTemplate(template);
    },
    updateInfraToEdit: async function(infraName){
      fetch("/api/infrastructures/"+infraName)
      .then(response => response.text())
      .then(text => { this.infraToEdit = text });
    },
    sendInfraToEdit: function(infraConfig){
      var value = document.getElementById('message-text-'+infraConfig).value;
      fetch('/api/infrastructures/'+infraConfig,{method: 'POST', body: value});
    },
    sendNewInfra: function(){
      var value = document.getElementById('textNewInfra').value;
      document.getElementById('textNewInfra').value = "";
      fetch('/api/infrastructures/'+value,{method: 'PUT'});
      location.reload();
    },
    deleteInfra: function(infraConfig){
      if (confirm('Are you sure you want to delete ' + infraConfig + ' test?')) {
        fetch('/api/infrastructures/'+infraConfig,{method: 'DELETE'});
        location.reload();
      }     
    },
    refresh: function() {
      if (this.infrastructureToLoad != ''){
        this.searchRunning();
        this.loadInfrasctructure(this.infrastructureToLoad);
      }
      setTimeout(this.refresh, 1000);
    },
    loadInfrasctructure(infrastructure){
      if (infrastructure == ''){
        return;
      }
      
      var runningIndex = 0;
      var runningFound = false;
      for (const runningInfra in this.running) {
        if (this.running[runningInfra].Name == infrastructure.Name){
          runningFound = true;
          break;
        }
        runningIndex += 1;
      }

      // Contenidor General del Test
      document.getElementById('infrastructure-zone').innerHTML = `
        <div class="container" id="container-infra-zone">
        </div>
      `

      // Capçalera del test
      document.getElementById('container-infra-zone').innerHTML = `
        <div class="row" id="infrastruture-row">
          <div class="col col-md">
            <div class="card">
              <div class="card-header">
                <b>Test Name:</b> ` + infrastructure.Name + `
                <button id="upinfrabutton" onClick="callApiUp('`+infrastructure.Name+`')" disabled>Up</button>
                <button id="runinfrabutton" onClick="callApiRun('`+infrastructure.Name+`')" disabled>Run</button>
                <button id="downinfrabutton" onClick="callApiDown('`+infrastructure.Name+`')" disabled>Down</button>
                <button id="cleaninfrabutton" onClick="callApiClean('`+infrastructure.Name+`')" disabled>Clean</button>
              </div>
              <div class="card-body">
                <p class="card-text">` + infrastructure.Description + `</p>
              </div>
            </div>
          </div>
        </div>
      `

      document.getElementById('personalizedScript').innerHTML = `
          function callApiUp(infra){
            fetch('/api/up/'+infra,{method: 'POST'});
          }
          function callApiDown(infra){
            fetch('/api/down/'+infra,{method: 'POST'});
          }
          function callApiRun(infra){
            fetch('/api/run/'+infra,{method: 'POST'})
          }
          function callApiClean(infra){
            fetch('/api/clean/'+infra,{method: 'POST'})
          }
          function showTestModal(infraName, hostI, testI){
            fetch('/api/infrastructures/'+infraName+'/hosts/'+hostI+'/tests/'+testI+'/result')
            .then(response => response.text())
            .then(text => { document.getElementById('modal-show-test-results-content').innerHTML = text });           
          }
      `
      
      if (runningFound){
        if (this.running[runningIndex].Status == "upping") {
          this.currentStatus = 'upping'
          document.getElementById('upinfrabutton').disabled = true;
          document.getElementById('runinfrabutton').disabled = true;
          document.getElementById('downinfrabutton').disabled = false;
          document.getElementById('cleaninfrabutton').disabled = true;
        }
        else if (this.running[runningIndex].Status == "up") {
          this.currentStatus = 'up'
          document.getElementById('upinfrabutton').disabled = true;
          document.getElementById('runinfrabutton').disabled = false;
          document.getElementById('downinfrabutton').disabled = false;
          document.getElementById('cleaninfrabutton').disabled = true;
        } 
        else if (this.running[runningIndex].Status == "running") {
          this.currentStatus = 'running'
          document.getElementById('upinfrabutton').disabled = true;
          document.getElementById('runinfrabutton').disabled = true;
          document.getElementById('downinfrabutton').disabled = false;
          document.getElementById('cleaninfrabutton').disabled = true;
        } 
        else if (this.running[runningIndex].Status == "done") {
          this.currentStatus = 'done'
          document.getElementById('upinfrabutton').disabled = true;
          document.getElementById('runinfrabutton').disabled = true;
          document.getElementById('downinfrabutton').disabled = false;
          document.getElementById('cleaninfrabutton').disabled = true;
        }
        else if (this.running[runningIndex].Status == "downing"){
          this.currentStatus = 'downing'
          this.refreshTime = '1000'
          document.getElementById('upinfrabutton').disabled = true;
          document.getElementById('runinfrabutton').disabled = true;
          document.getElementById('downinfrabutton').disabled = true;
          document.getElementById('cleaninfrabutton').disabled = true;
        } 
        else if (this.running[runningIndex].Status == "down") {
          this.currentStatus = 'down'
          document.getElementById('upinfrabutton').disabled = false;
          document.getElementById('runinfrabutton').disabled = true;
          document.getElementById('downinfrabutton').disabled = true;
          document.getElementById('cleaninfrabutton').disabled = false;
        }
        else if (this.running[runningIndex].Status == "cleaning") {
          this.currentStatus = 'cleaning'
          document.getElementById('upinfrabutton').disabled = true;
          document.getElementById('runinfrabutton').disabled = true;
          document.getElementById('downinfrabutton').disabled = true;
          document.getElementById('cleaninfrabutton').disabled = true;
        } 
      } else {
        this.currentStatus = ''
        document.getElementById('upinfrabutton').disabled = false;
        document.getElementById('runinfrabutton').disabled = true;
        document.getElementById('downinfrabutton').disabled = true;
        document.getElementById('cleaninfrabutton').disabled = true;
      }

      // Llistat de les xarxes
      for (const Net in infrastructure.Nets) {
        document.getElementById('container-infra-zone').innerHTML += `
          <div class="row" id="networks-row">
            <div class="col col-md">
              <div class="card">
                <div class="card-header"><b>Net Name:</b> `+ infrastructure.Nets[Net].Name +`</div>
                <div class="card-body">
                  <p class="card-text">IP: ` + infrastructure.Nets[Net].IP + `</p>
                </div>
              </div>
            </div>
          </div>
        `
      }

      // Contenidors pels hosts
      document.getElementById('container-infra-zone').innerHTML += `
        <div class="row" id="container-infra-zone-hosts">
        </div>`

      var nHost = 0;
      // Llistat de hosts
      for (const Host in infrastructure.Hosts) {

        var hostStyle = 'style="background: #fff"';
        var imageLoading = '';

        if (this.currentStatus == 'upping') {
          hostStyle = 'style="background: #EBF5FB"';
          imageLoading = '<img src="static/loading.gif" width=16/>';
        }
        else if (this.currentStatus == 'down'){
          hostStyle = 'style="background: #EBF5FB"';
          imageLoading = '';
        }
        else if (this.currentStatus == 'up' || this.currentStatus == 'running' || this.currentStatus == 'done' || this.currentStatus == 'downing') {
          hostStyle = 'style="background: #EAFAF1"';
          imageLoading = '';
        }

        if (runningFound && this.currentStatus == 'upping' && this.running[runningIndex].Hosts[nHost].Tests[0].Status == 1){
          hostStyle = 'style="background: #EAFAF1"';
          imageLoading = '';
        }

        document.getElementById('container-infra-zone-hosts').innerHTML += `
        <div class="col col-md" id="container-infra-zone-hosts-` + infrastructure.Hosts[Host].Name + `">
          <div class="card" id="card-host" ` + hostStyle + `>
            <div class="card-header">` + imageLoading + `<b> Host Name:</b> ` +  infrastructure.Hosts[Host].Name + `</div>
            <div class="card-body" id="container-infra-zone-hosts-` + infrastructure.Hosts[Host].Name + `-nets">
            <p class="card-text">Template: ` + infrastructure.Hosts[Host].Template + `</p>
            </div>
          </div>
          <div class="card" id="card-host-test-tile">
            <div class="card-header""><b>TESTS</b></div>
          </div>
        </div>
        `

        for (const Net in infrastructure.Hosts[Host].Nets) {
          document.getElementById('container-infra-zone-hosts-' + infrastructure.Hosts[Host].Name + '-nets').innerHTML += `
            <hr>
            <p class="card-text">Net: ` + infrastructure.Hosts[Host].Nets[Net].Name + `</p>
            <p class="card-text">IP: ` + infrastructure.Hosts[Host].Nets[Net].IP + `</p>
          `
        }

        //Llistat de tests
        var nTest = 1;
        var antTestStatus = 1;
        for (const Test in infrastructure.Hosts[Host].Tests) {
          var dependencies = ""
          if ( infrastructure.Hosts[Host].Tests[Test].Dependencies == null ) {
            dependencies = ""
          } else {
            dependencies = '<p class="card-text"> <b>Dep: </b>' + infrastructure.Hosts[Host].Tests[Test].Dependencies[0].Host + " - " + infrastructure.Hosts[Host].Tests[Test].Dependencies[0].Test + '</p>'
          }

          var args = ""
          if ( infrastructure.Hosts[Host].Tests[Test].Args == null ) {
            args = ""
          } else {
            args = "args: " + infrastructure.Hosts[Host].Tests[Test].Args
          }

          var testType = ""
          if ( infrastructure.Hosts[Host].Tests[Test].TestType == 0 ) {
            testType = "TestStopIfError"
          }
          else if ( infrastructure.Hosts[Host].Tests[Test].TestType == 1 ){
            testType = "TestOutputOnly"
          } 
          else if ( infrastructure.Hosts[Host].Tests[Test].TestType == 2 ){         
            testType = "TestDontStop"
          }

          var modeType = ""
          if ( infrastructure.Hosts[Host].Tests[Test].Mode == 0 ) {
            modeType = "NormalMode"
          }
          else if ( infrastructure.Hosts[Host].Tests[Test].Mode == 1 ){
            modeType = "RescuMode"
          } 

          var testStyle = ""
          if (runningFound){
            let testStatus = this.running[runningIndex].Hosts[nHost].Tests[nTest].Status;
            if (testStatus == 0 ) {
              if (antTestStatus == 1 && this.currentStatus == 'running'){
                //testStyle = 'style="background: #fff"'
                testStyle = 'style="background: #EBF5FB"';
                imageLoading = '<img src="static/loading.gif" width=16/>';
              } 
              else {
                testStyle = 'style="background: #fff"'
                imageLoading = ''
              }
            }
            else if (testStatus == 1 ) {
              testStyle = 'style="background: #EAFAF1"'
              imageLoading = ''
            }  
            else if (testStatus == -1 ) {
              testStyle = 'style="background: #F9EBEA"'
              imageLoading = ''
            }  
            antTestStatus = testStatus
          }      

          document.getElementById('container-infra-zone-hosts-' + infrastructure.Hosts[Host].Name).innerHTML += `
            <div class="col col-md" id="container-infra-zone-hosts-` + infrastructure.Hosts[Host].Name + `-` + nTest + `">
              <div class="card" data-bs-toggle="modal" data-bs-target="#modal-show-test-results" id="card-test" ` + testStyle + ` onclick="showTestModal('`+infrastructure.Name+`',` + nHost + `,` + nTest + `)">
                <div class="card-header">` + imageLoading + ` ` + nTest + ` - ` + infrastructure.Hosts[Host].Tests[Test].Name + ` </div>
                <div class="card-body">
                  <p class="card-text"> <b>Type: </b>` + testType + `</p>
                  <p class="card-text"> <b>Mode: </b>` + modeType + `</p>
                  `+ dependencies +`
                </div>
              </div>
            </div>
          `
          nTest += 1;
        }

        nHost += 1;
      }
    },
    loadTemplate(template){
      if (template == ''){
        return;
      }

      // Contenidor General dels Templates
      document.getElementById('infrastructure-zone').innerHTML = `
      <div class="accordion" id="accordionTTests">
        </div>
      `
      // Llistat del template tests
      for (const TTest in template.TemplateTests) {
        document.getElementById('accordionTTests').innerHTML += `

        <div class="accordion-item">
          <h5 class="accordion-header" id="heading`+ template.TemplateTests[TTest].Name.split('.').slice(0, -1).join('.') +`">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapse`+ template.TemplateTests[TTest].Name.split('.').slice(0, -1).join('.') +`" aria-expanded="false" aria-controls="collapse`+ template.TemplateTests[TTest].Name.split('.').slice(0, -1).join('.') +`">
              `+ template.TemplateTests[TTest].Name +`
            </button>
          </h5>
          <div id="collapse`+ template.TemplateTests[TTest].Name.split('.').slice(0, -1).join('.') +`" class="accordion-collapse collapse" aria-labelledby="heading`+ template.TemplateTests[TTest].Name.split('.').slice(0, -1).join('.') +`" data-bs-parent="#accordionTTests">
            <div class="accordion-body" style="white-space: pre-line">
              ` + template.TemplateTests[TTest].Content + `
            </div>
          </div>
        </div>
        `
      }
    }
  }
};

export default InfrastructuresTable;