var AppFooter = {
  props: ['context'],
  template: `
    <div v-if="context.loggedin" class="row">
      <br><br><br><br><br><br><br>
      <p style="text-align: center;">2021 - DRLM v2 Testing API</p>
    </div>
  `
}

export default AppFooter;
