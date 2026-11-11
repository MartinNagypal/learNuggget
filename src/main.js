import { mount } from 'svelte'
import App from './App.svelte'

import '@fontsource/nunito/800.css';
import '@fontsource/atkinson-hyperlegible/400.css';
import '@fontsource/atkinson-hyperlegible/700.css';

import '@fortawesome/fontawesome-free/css/fontawesome.min.css';
import '@fortawesome/fontawesome-free/css/solid.min.css';
import '@fortawesome/fontawesome-free/css/regular.min.css';

import './style/main.css';

const app = mount(App, {
  target: document.getElementById('app'),
})

export default app
