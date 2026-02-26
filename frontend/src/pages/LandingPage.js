import React from 'react';
import Navbar from '../components/Navbar';
import HeroSection from '../components/HeroSection';
import FeaturesSection from '../components/FeaturesSection';
import Footer from '../components/Footer';

import HowItWorksSection from '../components/HowItWorksSection';
import UseCasesSection from '../components/UseCasesSection';

function LandingPage() {
  return (
    <>
      <Navbar />
      <HeroSection />
      <FeaturesSection />
      <HowItWorksSection />
      <UseCasesSection />
      <Footer />
    </>
  );
}

export default LandingPage;
