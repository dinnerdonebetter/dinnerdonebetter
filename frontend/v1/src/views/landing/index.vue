<template>
  <div class="landing-container">
    <!-- Header element -->
    <header>
      <img src="@/assets/images/landing/header_logo.png">
      <h1>
        <a href="https://sacarter.design">Our Story</a> What's New? Recipe Book Cooking Class Sign Up
      </h1>
    </header>

    <!-- Main element -->
    <main>
      <h2 class="caption">Allons-y!</h2>

      <p>The community behind your cooking.</p>
      <img src="@/assets/images/landing/splash.png">
      <p class="caption">
        The <strong><em>place</em>&nbsp;</strong>for your <strong><em>mise.</em></strong>
      </p>

      <!-- Features -->
      <h2>App Features</h2>
      <ol>
        <li>Culinary Skill Development</li>
        <li>Recipe Customization</li>
        <li>Community Based Cooking</li>
      </ol>

      <!-- Available Platforms -->
      <h2>Download PrixFixe Today!</h2>
      <p>Available on iOS and Android:</p>
      <ul>
        <li>
          <a target="_blank" href="https://www.apple.com/ios/app-store/">Apple App Store</a>
        </li>
        <li>
          <a target="_blank" href="https://play.google.com/store?hl=en_US">Google Play Store</a>
        </li>
      </ul>
    </main>

    <!-- Footer Element -->
    <footer>
      <p>
        Design and Development by
        <a href="https://sacarter.design">Scott</a> &
        <a href="https://blog.verygoodsoftwarenotvirus.ru">Jeffrey</a>, respectively.
      </p>
    </footer>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component({
  name: 'Landing',
})
export default class extends Vue {
  mounted() {
    document.title = "PrixFixe - the place for your mise";
  }
}
</script>

<style lang="scss" scoped>
@font-face {
  font-family: 'MADE TOMMY';
  src: url(~@/assets/fonts/MadeTommy/madetommy_regular_macroman/made_tommy_regular-webfont.woff2) format('woff2');
}

body {
  background-color: #f5f6ff;
  margin: auto;
  font-family: MADE TOMMY, Helvetica, Arial, sans-serif;
  font-size: 20px;
  line-height: 30px;
}

/* header styles */

h1,
h2,
p,
ol,
li {
  color: #232d61;
}

h1,
h2 {
  font-size: 41px;
}

p,
ol,
li {
  font-size: 25px;
}

header {
  height: 128px;
  width: 1920px;
  background-color: #ffffff;
}

header h1 {
  color: #232d61;
  font-family: MADE TOMMY, Helvetica, Arial, sans-serif;
  font-size: 25px;
  min-height: 4em;
  text-align: center;
  display: inline-block;
  vertical-align: middle;
}

header img {
  padding: 0 100px;
  padding-top: 25px;
}

/*main body styles*/

main {
  max-width: 1920px;
  margin: 0 auto;
  padding: 25px 100px;
}

main img {
  max-width: 100%;
}

a {
  color: #dc5126;
}
a:hover {
  color: #232d61;
}

/*footer styles*/

footer {
  height: 32px;
  width: 1920px;
  background-color: #eeeffe;
  text-align: center;
  padding: 50px;
}

footer p {
  color: #232d61;
  font-family: MADE TOMMY, Helvetica, Arial, sans-serif;
  font-size: 25px;
  min-height: 4em;
  text-align: center;
  display: inline-block;
  vertical-align: middle;
  margin: 0;
}

p.caption {
  color: #dc5126;
}

h2.caption {
  color: #dc5126;
}
</style>
