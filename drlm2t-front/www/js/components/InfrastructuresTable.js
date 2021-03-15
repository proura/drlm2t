var InfrastructuresTable = {
  props: ['context'],
  data() {
    return {
      infrastructures: 'loading',
      running: 'loading',
      templates: 'loading',
      infrastructureToLoad: '',
      templateToLoad: '',
      infraToEdit: 'test'
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
                        <form>
                          <div class="mb-3">
                            <textarea class="form-control" :bind="infraToEdit" :id="'message-text-'+ infrastructure.Name">{{ infraToEdit }}</textarea>
                          </div>
                        </form>
                      </div>
                      <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                        <button type="button" class="btn btn-primary" v-on:click="sendInfraToEdit(infrastructure.Name)">Save</button>
                      </div>
                    </div>
                  </div>
                </div>
              </tr>
            </tbody>
          </table>
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
    </div>
  `,
  created() {
    this.searchInfrastructure();
    this.refresh();
    //this.searchRunning();
    this.searchTemplates();
  },
  methods: {
    searchInfrastructure: function () {
      fetch('/api/infrastructures')
        .then(response => response.json())
        .then(json => {this.infrastructures = json.resultList.result;});
    },
    searchRunning: function () {
      console.log("fetch running!!!")
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
      //this.templateToLoad = template;
      this.loadTemplate(template)
    },
    updateInfraToEdit: async function(infraName){
      fetch("/api/infrastructures/"+infraName)
      .then(response => response.text())
      .then(text => { this.infraToEdit = text });
    },
    sendInfraToEdit: function(infraConfig){
      var value = document.getElementById('message-text-'+infraConfig).value;

      console.log(value)
      fetch('/api/infrastructures/'+infraConfig,{method: 'POST', body: value})
      // .then(response => {
      //   if (response.ok){
      //     this.context.loggedin = false;
      //   }
      // });
    },
    refresh: function() {
      if (this.infrastructureToLoad != ''){
        this.searchRunning()
        this.loadInfrasctructure(this.infrastructureToLoad)
      }
      setTimeout(this.refresh,3000)
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
                </form id="upinfrabutton"><button onClick="callApiUp('`+infrastructure.Name+`')">Up</button></form>
                </form id="runinfrabutton"><button onClick="callApiRun('`+infrastructure.Name+`')">Run</button></form>
                </form id="downinfrabutton"><button onClick="callApiDown('`+infrastructure.Name+`')">Down</button></form>
                </form id="cleaninfrabutton"><button onClick="callApiClean('`+infrastructure.Name+`')">Clean</button></form>
              </div>
              <div class="card-body">
                <p class="card-text">` + infrastructure.Description + `</p>
              </div>
            </div>
          </div>
        </div>
      `

      var script = document.createElement("script");
      script.innerHTML = `
        function callApiUp(infra){
          console.log("UP: " + infra);
          fetch('/api/up/'+infra,{method: 'POST'});
          //.then(response => {if (response.ok){this.context.loggedin = false;}});
          //event.preventDefault();
        }
        function callApiDown(infra){
          console.log("DOWN: " + infra);
          fetch('/api/down/'+infra,{method: 'POST'});
          //.then(response => {if (response.ok){this.context.loggedin = false;}});
          //event.preventDefault();
        }
        function callApiRun(infra){
          console.log("RUN: " + infra);
          fetch('/api/run/'+infra,{method: 'POST'})
          //.then(response => {if (response.ok){this.context.loggedin = false;}});
          //event.preventDefault();
        }
        function callApiClean(infra){
          console.log("CLEAN: " + infra);
          fetch('/api/clean/'+infra,{method: 'POST'})
          //.then(response => {if (response.ok){this.context.loggedin = false;}});
          //event.preventDefault();
        }
      `
      document.body.appendChild(script);

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

        document.getElementById('container-infra-zone-hosts').innerHTML += `
        <div class="col col-md" id="container-infra-zone-hosts-` + infrastructure.Hosts[Host].Name + `">
          <div class="card" id="card-host">
            <div class="card-header"><b>Host Name:</b> ` +  infrastructure.Hosts[Host].Name + `</div>
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

          var style = ""
          if (runningFound){
            let testStatus = this.running[runningIndex].Hosts[nHost].Tests[nTest].Status;
            if (testStatus == 0 ) {
              style = 'style="background: #fff"'
            }
            else if (testStatus == 1 ) {
              style = 'style="background: #b1e19d"'
            }  
            else if (testStatus == 2 ) {
              style = 'style="background: #ff9d9d"'
            }  
          }
                    
          document.getElementById('container-infra-zone-hosts-' + infrastructure.Hosts[Host].Name).innerHTML += `
            <div class="col col-md" id="container-infra-zone-hosts-` + infrastructure.Hosts[Host].Name + `-` + nTest + `">
              <div class="card" id="card-test" ` + style + `>
                <div class="card-header">` + nTest + ` - ` + infrastructure.Hosts[Host].Tests[Test].Name + `</div>
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