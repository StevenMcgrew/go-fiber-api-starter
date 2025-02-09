import NotFoundPage from './lib/pages/NotFoundPage.svelte';
import HomePage from './lib/pages/HomePage.svelte';
import AboutPage from './lib/pages/AboutPage.svelte';
import ContactPage from './lib/pages/ContactPage.svelte';
import FailedEmailVerificationPage from './lib/pages/FailedEmailVerificationPage.svelte';
import SuccessEmailVerificationPage from './lib/pages/SuccessEmailVerificationPage.svelte';

export const routes = {
  '/': HomePage,
  '/about': AboutPage,
  '/contact': ContactPage,
  '/failed-email-verification': FailedEmailVerificationPage,
  '/success-email-verification': SuccessEmailVerificationPage,
  '*': NotFoundPage
};