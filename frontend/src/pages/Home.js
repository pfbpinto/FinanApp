import React from "react";

function Home() {
  return (
    <div className="flex flex-col items-center justify-center px-4 sm:px-6 lg:px-8">
      {/* Conteúdo principal */}
      <div className="text-center max-w-3xl py-12">
        <h1 className="text-4xl font-extrabold text-gray-900 sm:text-5xl">
          Welcome to FinanAPP
        </h1>
        <p className="mt-4 text-lg text-gray-600">
          A great place to find awesome content, services, and products.
        </p>
      </div>

      {/* Botões de Ação */}
      <div className="mt-8 flex justify-center space-x-6">
        <a
          href="#features"
          className="inline-block rounded-md bg-blue-600 text-white px-6 py-3 text-lg font-semibold hover:bg-blue-700"
        >
          Explore Features
        </a>
        <a
          href="#contact"
          className="inline-block rounded-md bg-transparent border-2 border-blue-600 text-blue-600 px-6 py-3 text-lg font-semibold hover:bg-blue-100"
        >
          Contact Us
        </a>
      </div>

      {/* Recursos */}
      <div
        id="features"
        className="mt-16 grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl"
      >
        <div className="text-center bg-white p-6 rounded-lg shadow-lg">
          <h3 className="text-2xl font-semibold text-gray-900">Feature One</h3>
          <p className="mt-4 text-gray-600">
            Description of the first amazing feature goes here. It's a simple
            and effective feature.
          </p>
        </div>
        <div className="text-center bg-white p-6 rounded-lg shadow-lg">
          <h3 className="text-2xl font-semibold text-gray-900">Feature Two</h3>
          <p className="mt-4 text-gray-600">
            Description of the second amazing feature goes here. It offers great
            value to users.
          </p>
        </div>
      </div>

      {/* Contato */}
      <div id="contact" className="mt-16 bg-blue-600 text-white py-12 w-full">
        <div className="text-center">
          <h2 className="text-3xl font-semibold">Get in Touch</h2>
          <p className="mt-4 text-lg">Have questions? We're here to help!</p>
          <a
            href="mailto:info@website.com"
            className="mt-8 inline-block rounded-md bg-white text-blue-600 px-6 py-3 text-lg font-semibold hover:bg-gray-100"
          >
            Contact Us via Email
          </a>
        </div>
      </div>
    </div>
  );
}

export default Home;
