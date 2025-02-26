import NotFoundPage from './lib/pages/NotFoundPage.svelte';
import HomePage from './lib/pages/HomePage.svelte';
import AboutPage from './lib/pages/AboutPage.svelte';
import ContactPage from './lib/pages/ContactPage.svelte';
import AccountPage from './lib/pages/AccountPage.svelte';

export const routes = {
  '/': HomePage,
  '/about': AboutPage,
  '/contact': ContactPage,
  '/account': AccountPage,
  '*': NotFoundPage,
};