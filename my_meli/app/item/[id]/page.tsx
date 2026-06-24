'use client';

import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import { useParams } from 'next/navigation';

export default function ItemPage() {
  const params = useParams();
  const itemId = (params?.id as string) || "MLA43960787";

  const [itemData, setItemData] = useState<{ title: string, price: number } | null>(null);
  const [fetchError, setFetchError] = useState<string | null>(null);
  // Interactive states
  const [cartCount, setCartCount] = useState(3);
  const [selectedQty, setSelectedQty] = useState(1);
  const [activeTab, setActiveTab] = useState(0); // 0 = Imagen Principal
  const [showCheckout, setShowCheckout] = useState(false);
  const [loadingStep, setLoadingStep] = useState(0); // 0 = cerrado, 1 = cargando, 2 = éxito
  const [gatewayResponse, setGatewayResponse] = useState<any>(null);
  const [showToast, setShowToast] = useState(false);
  const [selectedPayment, setSelectedPayment] = useState('tarjeta-mp');

  // Dynamic product data or fallback
  const productPrice = itemData?.price || 0;
  const originalPrice = productPrice > 0 ? productPrice * 1.25 : 0;
  const productName = itemData?.title || (fetchError ? "Error cargando producto" : "Cargando producto...");

  useEffect(() => {
    if (!itemId) return;
    fetch(`http://localhost:8080/gateway/items/${itemId}`)
      .then(res => res.json())
      .then(data => {
        if (data.error) {
          setFetchError(data.message || data.error);
        } else if (data.data) {
          setItemData({ title: data.data.title, price: data.data.price });
        } else {
          setFetchError("Respuesta inesperada del servidor");
        }
      })
      .catch(err => setFetchError("No se pudo contactar al Gateway para obtener el item."));
  }, [itemId]);

  // E2E purchase event simulation to API Gateway
  const handleConfirmPurchase = async () => {
    setLoadingStep(1); // Inicia pantalla de carga

    // Simulate processing delay
    await new Promise(resolve => setTimeout(resolve, 2000));

    const purchasePayload = {
      event: "purchase_initiated",
      user: "Eliseo",
      item_id: itemId,
      quantity: selectedQty,
      address: "Calle Falsa 123",
      amount: productPrice * selectedQty,
      timestamp: new Date().toISOString()
    };

    try {
      // 1. Authenticate to get the JWT token
      const authResponse = await fetch('http://localhost:8080/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username: 'admin', password: 'admin' })
      });

      if (!authResponse.ok) {
        setGatewayResponse({ error: "Auth Failed", message: "No se pudo obtener el token JWT." });
        setLoadingStep(2);
        return;
      }

      const authData = await authResponse.json();
      const token = authData.data.token;

      // 2. Send event to Nginx port (8080) with the JWT
      const response = await fetch('http://localhost:8080/gateway/orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(purchasePayload)
      });

      if (response.ok) {
        const data = await response.json();
        setGatewayResponse(data);
      } else {
        setGatewayResponse({ error: `Error HTTP: ${response.status}`, message: "Acceso denegado o error interno" });
      }
    } catch (error) {
      console.warn("No se pudo conectar con el Gateway en el puerto 8080. ¿Están corriendo Nginx y FastAPI?", error);
      setGatewayResponse({
        error: "Gateway Offline",
        message: "No se pudo contactar al balanceador (http://localhost:8080/gateway). Asegúrate de tener los servicios activos."
      });
    }

    setLoadingStep(2); // Muestra pantalla de éxito
  };

  // Add to cart handler
  const handleAddToCart = () => {
    setCartCount(prev => prev + selectedQty);
    setShowToast(true);
    setTimeout(() => setShowToast(false), 3000);
  };

  return (
    <div className="min-h-screen bg-gray-100 font-sans text-gray-800 antialiased pb-12">
      {/* HEADER DE MERCADOLIBRE (ALTA FIDELIDAD) */}
      <header className="bg-[#fff159] border-b border-gray-300 py-2 sticky top-0 z-40 shadow-sm">
        <div className="max-w-6xl mx-auto px-4">
          {/* Fila Superior */}
          <div className="flex items-center justify-between gap-4">
            {/* Logo */}
            <div className="flex items-center gap-2 cursor-pointer">
              <div className="w-11 h-8 bg-black rounded flex items-center justify-center font-bold text-white text-xs tracking-tighter">
                MELI
              </div>
              <span className="font-extrabold text-blue-900 text-lg leading-none hidden sm:inline">mercado<br /><span className="text-gray-900 font-bold">libre</span></span>
            </div>

            {/* Barra de Búsqueda */}
            <div className="flex-1 max-w-2xl relative">
              <input
                type="text"
                placeholder="Buscar productos, marcas y más..."
                className="w-full bg-white text-gray-800 pl-4 pr-10 py-2 rounded-sm shadow-sm text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 placeholder-gray-400"
                defaultValue="Monitor gamer curvo Xiaomi Gaming G34WQi"
              />
              <button className="absolute right-3 top-2.5 text-gray-400">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24"><path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
              </button>
            </div>

            {/* Banner Meli+ */}
            <div className="hidden md:flex items-center bg-[#181824] text-white text-xs px-3 py-1.5 rounded-sm gap-2 cursor-pointer hover:bg-opacity-90">
              <span className="font-bold text-[#fff159]">meli+</span>
              <span className="opacity-80">POR $ 20.990</span>
              <span className="font-bold text-green-400">$ 11.999/MES</span>
            </div>
          </div>

          {/* Fila Inferior */}
          <div className="flex items-center justify-between mt-3 text-xs text-gray-700">
            {/* Dirección de Envío */}
            <div className="flex items-center gap-1.5 cursor-pointer group hover:bg-yellow-200 px-2 py-1 rounded">
              <svg className="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"></path><path strokeLinecap="round" strokeLinejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"></path></svg>
              <div className="leading-tight">
                <span className="text-[10px] text-gray-500 block">Enviar a Eliseo</span>
                <span className="font-semibold block text-gray-900">Calle Street Number</span>
              </div>
            </div>

            {/* Enlaces de Navegación */}
            <nav className="hidden lg:flex items-center gap-4 text-gray-600">
              <span className="cursor-pointer hover:text-black">Categorías</span>
              <span className="cursor-pointer hover:text-black">Ofertas</span>
              <span className="cursor-pointer hover:text-black">Cupones</span>
              <span className="cursor-pointer hover:text-black">Supermercado</span>
              <span className="cursor-pointer hover:text-black">Moda</span>
              <span className="cursor-pointer hover:text-black">Mercado Play</span>
              <span className="cursor-pointer hover:text-black">Vender</span>
              <span className="cursor-pointer hover:text-black">Ayuda</span>
            </nav>

            {/* Perfil y Acciones de Eliseo */}
            <div className="flex items-center gap-4">
              {/* Usuario */}
              <div className="flex items-center gap-1 cursor-pointer">
                <div className="w-5 h-5 bg-gray-400 rounded-full flex items-center justify-center text-[10px] font-bold text-white">E</div>
                <span className="font-semibold text-gray-900">Eliseo</span>
                <svg className="w-2.5 h-2.5 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"></path></svg>
              </div>

              {/* Menú */}
              <span className="cursor-pointer hover:text-black hidden sm:inline">Mis compras</span>
              <div className="flex items-center gap-1 cursor-pointer hover:text-black hidden sm:flex">
                <span>Favoritos</span>
                <svg className="w-2.5 h-2.5 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"></path></svg>
              </div>

              {/* Notificaciones */}
              <div className="relative cursor-pointer text-gray-600 hover:text-black">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0"></path></svg>
                <span className="absolute -top-1.5 -right-1.5 bg-red-600 text-white rounded-full text-[9px] font-bold px-1 py-0.2">3</span>
              </div>

              {/* Carrito */}
              <div className="relative cursor-pointer text-gray-600 hover:text-black">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 00-16.536-1.84M7.5 14.25L5.106 5.272M6 20.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm12.75 0a.75.75 0 11-1.5 0 .75.75 0 011.5 0z"></path></svg>
                {cartCount > 0 && (
                  <span className="absolute -top-1.5 -right-1.5 bg-red-600 text-white rounded-full text-[9px] font-bold px-1 py-0.2">{cartCount}</span>
                )}
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* CUERPO PRINCIPAL DEL PRODUCTO */}
      <main className="max-w-6xl mx-auto px-4 mt-6">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-xs text-blue-600 mb-4">
          <span className="cursor-pointer font-semibold hover:underline">Volver</span>
          <span className="text-gray-400">|</span>
          <span className="cursor-pointer hover:underline text-gray-500">Computación</span>
          <span className="text-gray-400">&gt;</span>
          <span className="cursor-pointer hover:underline text-gray-500">Monitores y Accesorios</span>
          <span className="text-gray-400">&gt;</span>
          <span className="cursor-pointer hover:underline text-gray-500">Monitores</span>
        </div>

        {/* Grid de 2 Columnas de Mercado Libre */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">

          {/* COLUMNA IZQUIERDA (Imágenes y Ficha Técnica) */}
          <div className="lg:col-span-8 bg-white rounded-lg p-6 border border-gray-200 shadow-xs">
            <div className="flex flex-col md:flex-row gap-6">
              {/* Selector de Miniaturas */}
              <div className="flex md:flex-col gap-2 order-2 md:order-1 justify-center">
                {[0, 1, 2].map((idx) => (
                  <button
                    key={idx}
                    onClick={() => setActiveTab(idx)}
                    className={`w-12 h-12 border-2 rounded overflow-hidden p-0.5 flex items-center justify-center bg-gray-50 transition-all ${activeTab === idx ? 'border-blue-500' : 'border-gray-200 hover:border-gray-400'}`}
                  >
                    <Image
                      src="/xiaomi_monitor.png"
                      alt="Thumbnail"
                      width={48}
                      height={48}
                      className={`object-contain transition-transform ${idx === 1 ? 'rotate-6' : idx === 2 ? '-rotate-6' : ''}`}
                      style={{ width: "auto", height: "auto" }}
                    />
                  </button>
                ))}
              </div>

              {/* Imagen Principal Grande */}
              <div className="flex-1 flex justify-center items-center order-1 md:order-2 p-4 min-h-[350px] relative border border-gray-100 rounded-md">
                <Image
                  src="/xiaomi_monitor.png"
                  alt="Monitor gamer curvo Xiaomi Gaming G34WQi"
                  width={450}
                  height={450}
                  className={`object-contain max-h-[400px] transition-transform duration-300 ${activeTab === 1 ? 'scale-105' : activeTab === 2 ? 'rotate-2' : ''}`}
                  style={{ width: "auto", height: "auto" }}
                  priority
                />
              </div>
            </div>

            {/* Ficha Técnica (Specs) */}
            <div className="mt-10 border-t border-gray-200 pt-8">
              <h2 className="text-xl font-semibold mb-6">Lo que tenés que saber de este producto</h2>

              <ul className="grid grid-cols-1 md:grid-cols-2 gap-y-4 gap-x-6 text-sm text-gray-700">
                <li className="flex items-start gap-3">
                  <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0"></span>
                  <div><strong className="text-gray-900">Voltaje:</strong> 220V.</div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0"></span>
                  <div><strong className="text-gray-900">Tamaño de la pantalla:</strong> 34 pulgadas.</div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0"></span>
                  <div><strong className="text-gray-900">Diseño:</strong> Monitor Curvo.</div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0"></span>
                  <div><strong className="text-gray-900">Pantalla:</strong> Posee pantalla antirreflejo.</div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0"></span>
                  <div><strong className="text-gray-900">Resolución:</strong> UltraWide WQHD de 3440px-1440px.</div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0"></span>
                  <div><strong className="text-gray-900">Relación de aspecto:</strong> 21:9 para mayor inmersión.</div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0"></span>
                  <div><strong className="text-gray-900">Panel:</strong> Panel VA de alto contraste.</div>
                </li>
              </ul>
            </div>
          </div>

          {/* COLUMNA DERECHA (Caja de Compra y Datos de Vendedor) */}
          <div className="lg:col-span-4 flex flex-col gap-4">

            {/* Tarjeta de compra principal */}
            <div className="bg-white rounded-lg p-6 border border-gray-200 shadow-xs">
              <span className="text-xs text-gray-500 font-medium block">Nuevo | +100 vendidos</span>

              {/* Título y Favorito */}
              <div className="flex items-start justify-between gap-3 mt-1">
                <h1 className="text-lg font-bold text-gray-900 leading-tight">
                  {productName}
                </h1>
                <button className="text-blue-500 hover:scale-115 transition-transform shrink-0">
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"></path></svg>
                </button>
              </div>

              {/* Puntuación */}
              <div className="flex items-center gap-1.5 mt-2">
                <span className="text-xs font-semibold text-gray-800">4.8</span>
                <div className="flex text-blue-500">
                  {[1, 2, 3, 4, 5].map((s) => (
                    <svg key={s} className="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24"><path d="M12 .587l3.668 7.431 8.2 1.192-5.934 5.787 1.4 8.168L12 18.896l-7.334 3.857 1.4-8.168L.132 9.21l8.2-1.192L12 .587z"></path></svg>
                  ))}
                </div>
                <span className="text-xs text-gray-400 font-medium">(1227)</span>
              </div>

              {/* Precios */}
              <div className="mt-4">
                <span className="text-xs text-gray-400 line-through">$ {originalPrice.toLocaleString('es-AR')}</span>
                <div className="flex items-center gap-2 mt-0.5">
                  <span className="text-3xl font-normal text-gray-900">$ {productPrice.toLocaleString('es-AR')}</span>
                  <span className="text-sm font-bold text-green-500">20% OFF</span>
                </div>
                <span className="text-xs font-semibold text-green-600 block mt-1.5">
                  Mismo precio en 3 cuotas de $ {(productPrice / 3).toLocaleString('es-AR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} con tarjeta Mercado Pago
                </span>
                <span className="text-[10px] text-gray-400 block mt-1">Precio sin impuestos nacionales: $ 469.301</span>

                <span className="text-xs text-blue-500 hover:underline cursor-pointer font-medium block mt-3">
                  Ver los medios de pago
                </span>
              </div>

              {/* Color */}
              <div className="mt-6">
                <span className="text-xs text-gray-500 font-medium">Color: <strong className="text-gray-900">Negro</strong></span>
                <div className="w-10 h-10 border-2 border-blue-500 rounded p-0.5 mt-2 flex items-center justify-center bg-gray-50 cursor-pointer">
                  <Image src="/xiaomi_monitor.png" alt="Color Negro" width={32} height={32} className="object-contain" style={{ width: "auto", height: "auto" }} />
                </div>
              </div>

              {/* Envío */}
              <div className="mt-6 flex items-start gap-3 bg-gray-50 p-3 rounded-md border border-gray-100">
                <svg className="w-5 h-5 text-green-600 mt-0.5 shrink-0" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M8.25 18.75a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 01-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h1.125a1.125 1.125 0 001.125-1.125V9.75M8.25 18.75h7.5M12 15.75h.008v.008H12v-.008z"></path></svg>
                <div>
                  <span className="text-sm font-semibold text-green-600 block">Llega gratis a partir del lunes</span>
                  <span className="text-xs text-gray-500">A Calle Falsa 123, CABA</span>
                  <span className="text-xs text-blue-500 block hover:underline cursor-pointer font-medium mt-1">Más detalles y formas de entrega</span>
                </div>
              </div>

              {/* Cantidad Selector */}
              <div className="mt-5 text-sm">
                <span className="font-semibold block text-gray-900">Stock disponible</span>
                <div className="flex items-center gap-2 mt-2">
                  <label htmlFor="qty-select" className="text-xs text-gray-500">Cantidad:</label>
                  <select
                    id="qty-select"
                    value={selectedQty}
                    onChange={(e) => setSelectedQty(parseInt(e.target.value))}
                    className="bg-white border border-gray-300 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-blue-500 font-semibold text-sm cursor-pointer"
                  >
                    {[1, 2, 3, 4, 5].map((n) => (
                      <option key={n} value={n}>{n} unidad{n > 1 ? 'es' : ''}</option>
                    ))}
                  </select>
                  <span className="text-xs text-gray-400 font-medium">(+50 disponibles)</span>
                </div>
              </div>

              {/* Botones de Acción */}
              <div className="mt-6 flex flex-col gap-2">
                <button
                  onClick={() => setShowCheckout(true)}
                  className="w-full bg-[#3483fa] text-white py-3 rounded-md font-bold text-sm cursor-pointer shadow-sm hover:bg-opacity-95 transition-all text-center"
                >
                  Comprar ahora
                </button>
                <button
                  onClick={handleAddToCart}
                  className="w-full bg-[#e3edfb] text-[#3483fa] py-3 rounded-md font-bold text-sm cursor-pointer hover:bg-opacity-90 transition-all text-center"
                >
                  Agregar al carrito
                </button>
              </div>

              {/* Tienda Oficial */}
              <div className="mt-6 border-t border-gray-100 pt-4 text-xs text-gray-600 flex items-center justify-between">
                <div>
                  <span className="block text-gray-400">Tienda oficial</span>
                  <span className="font-semibold text-gray-900 block mt-0.5">Electro World Group</span>
                </div>
                <svg className="w-5 h-5 text-blue-500 fill-current" viewBox="0 0 24 24"><path d="M12 0c-6.627 0-12 5.373-12 12s5.373 12 12 12 12-5.373 12-12-5.373-12-12-12zm-1.25 17.292l-4.5-4.364 1.857-1.858 2.643 2.506 5.643-5.784 1.857 1.857-7.5 7.643z"></path></svg>
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* TOAST DE AGREGAR AL CARRITO */}
      {showToast && (
        <div className="fixed bottom-5 right-5 bg-gray-900 text-white rounded-lg px-4 py-3 flex items-center gap-3 shadow-lg border border-gray-700 animate-slide-in z-50">
          <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          <div className="text-xs">
            <span className="font-bold block">¡Agregado al carrito!</span>
            <span>Añadiste {selectedQty} monitor{selectedQty > 1 ? 'es' : ''} a tu bolsa.</span>
          </div>
        </div>
      )}

      {/* MODAL 1: CHECKOUT (COMPRAR AHORA) */}
      {showCheckout && loadingStep === 0 && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50 animate-fade-in">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden animate-scale-up">

            {/* Header Checkout */}
            <div className="bg-[#fff159] p-4 border-b border-gray-300 flex justify-between items-center">
              <h3 className="font-bold text-gray-900 text-sm">Resumen de tu compra</h3>
              <button
                onClick={() => setShowCheckout(false)}
                className="text-gray-700 hover:text-black font-extrabold text-sm"
              >
                ✕
              </button>
            </div>

            {/* Contenido Checkout */}
            <div className="p-5 flex flex-col gap-4">

              {/* Dirección de Envío */}
              <div className="border border-gray-200 rounded p-3 bg-gray-50 flex gap-3 items-start">
                <svg className="w-5 h-5 text-blue-500 shrink-0 mt-0.5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"></path><path strokeLinecap="round" strokeLinejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"></path></svg>
                <div>
                  <span className="text-xs font-semibold block text-gray-900">Dirección de Entrega:</span>
                  <span className="text-xs text-gray-600 block">Eliseo - Calle Falsa 123, CABA</span>
                  <span className="text-[10px] text-gray-400">Llega gratis el lunes</span>
                </div>
              </div>

              {/* Producto a Comprar */}
              <div className="flex gap-3 items-center border-b border-gray-100 pb-3">
                <div className="w-12 h-12 bg-gray-50 border border-gray-200 rounded p-1 flex items-center justify-center shrink-0">
                  <Image src="/xiaomi_monitor.png" alt="Xiaomi Monitor" width={40} height={40} className="object-contain" style={{ width: "auto", height: "auto" }} />
                </div>
                <div className="flex-1 min-w-0">
                  <span className="text-xs font-semibold text-gray-800 truncate block">{productName}</span>
                  <span className="text-xs text-gray-500">Cantidad: {selectedQty} unidad{selectedQty > 1 ? 'es' : ''}</span>
                </div>
                <span className="text-xs font-bold text-gray-900 shrink-0">$ {(productPrice * selectedQty).toLocaleString('es-AR')}</span>
              </div>

              {/* Selector de Método de Pago */}
              <div>
                <span className="text-xs font-bold block text-gray-800 mb-2">Selecciona un medio de pago:</span>
                <div className="flex flex-col gap-2">
                  <label className="flex items-center gap-3 border border-gray-200 rounded p-2.5 cursor-pointer hover:bg-gray-50 text-xs">
                    <input
                      type="radio"
                      name="payment-method"
                      value="tarjeta-mp"
                      checked={selectedPayment === 'tarjeta-mp'}
                      onChange={() => setSelectedPayment('tarjeta-mp')}
                      className="accent-blue-500"
                    />
                    <div className="flex-1">
                      <span className="font-semibold block text-gray-800">Tarjeta Mercado Pago (Visa Débito)</span>
                      <span className="text-[10px] text-green-600">Saldo disponible / Sin cargo de procesamiento</span>
                    </div>
                  </label>
                  <label className="flex items-center gap-3 border border-gray-200 rounded p-2.5 cursor-pointer hover:bg-gray-50 text-xs">
                    <input
                      type="radio"
                      name="payment-method"
                      value="tarjeta-credito"
                      checked={selectedPayment === 'tarjeta-credito'}
                      onChange={() => setSelectedPayment('tarjeta-credito')}
                      className="accent-blue-500"
                    />
                    <div className="flex-1">
                      <span className="font-semibold block text-gray-800">Tarjeta de Crédito Bancaria</span>
                      <span className="text-[10px] text-gray-500">3 cuotas sin interés</span>
                    </div>
                  </label>
                </div>
              </div>

              {/* Totales */}
              <div className="bg-gray-50 p-3 rounded border border-gray-100 flex flex-col gap-1.5 text-xs">
                <div className="flex justify-between">
                  <span className="text-gray-500">Producto</span>
                  <span>$ {(productPrice * selectedQty).toLocaleString('es-AR')}</span>
                </div>
                <div className="flex justify-between text-green-600">
                  <span>Envío</span>
                  <span className="font-semibold">Gratis</span>
                </div>
                <div className="flex justify-between border-t border-gray-200 pt-2 font-bold text-sm text-gray-900">
                  <span>Total</span>
                  <span>$ {(productPrice * selectedQty).toLocaleString('es-AR')}</span>
                </div>
              </div>
            </div>

            {/* Footer Checkout Acción */}
            <div className="p-4 bg-gray-50 border-t border-gray-200 flex items-center justify-end gap-2">
              <button
                onClick={() => setShowCheckout(false)}
                className="px-4 py-2 text-xs font-semibold text-gray-600 hover:text-black"
              >
                Cancelar
              </button>
              <button
                onClick={handleConfirmPurchase}
                className="px-6 py-2.5 bg-blue-500 hover:bg-blue-600 text-white rounded font-bold text-xs shadow-sm cursor-pointer"
              >
                Confirmar Compra
              </button>
            </div>

          </div>
        </div>
      )}

      {/* MODAL 2: CARGANDO COMPRA Y ENVIANDO AL GATEWAY */}
      {showCheckout && loadingStep === 1 && (
        <div className="fixed inset-0 bg-black bg-opacity-65 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-xs p-6 flex flex-col items-center justify-center text-center gap-4 animate-scale-up">

            {/* Spinner */}
            <div className="relative w-12 h-12">
              <div className="absolute inset-0 rounded-full border-4 border-blue-100"></div>
              <div className="absolute inset-0 rounded-full border-4 border-blue-500 border-t-transparent animate-spin"></div>
            </div>

            <div className="flex flex-col gap-1">
              <span className="font-bold text-gray-900 text-sm">Procesando pago...</span>
              <span className="text-[10px] text-gray-500">Contactando con Mercado Pago y despachando evento E2E al Gateway en http://localhost:8080/gateway...</span>
            </div>

          </div>
        </div>
      )}

      {/* MODAL 3: ÉXITO (COMPRA FINALIZADA) */}
      {showCheckout && loadingStep === 2 && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden animate-scale-up">

            {/* Header éxito */}
            <div className="bg-green-500 p-6 text-white text-center flex flex-col items-center gap-2">
              {/* Checkmark Icon */}
              <div className="w-12 h-12 bg-white rounded-full flex items-center justify-center text-green-500 text-2xl font-bold shadow-md">
                ✓
              </div>
              <h3 className="font-bold text-lg">¡Excelente compra, Eliseo!</h3>
              <p className="text-xs opacity-90">Tu compra ha sido procesada de manera exitosa.</p>
            </div>

            {/* Detalles de la compra */}
            <div className="p-5 flex flex-col gap-4 text-xs">

              {/* Ubicación de envío y llegada */}
              <div className="border border-green-100 bg-green-50/50 rounded p-3">
                <span className="font-bold text-green-800 block text-[11px] mb-1">📅 Información de entrega:</span>
                <span className="text-gray-700 block font-medium">Llegará el lunes a tu dirección:</span>
                <span className="text-gray-900 font-bold block mt-0.5">Calle Falsa 123, CABA</span>
              </div>

              {/* Registro del balanceador de carga / Gateway */}
              <div className="border border-gray-200 rounded p-3 bg-gray-50">
                <span className="font-bold text-gray-800 block text-[11px] mb-1.5">📡 Diagnóstico del Gateway (E2E):</span>

                {gatewayResponse ? (
                  gatewayResponse.error ? (
                    <div className="text-yellow-600 bg-yellow-50 p-2 rounded border border-yellow-200">
                      <span className="font-bold block">⚠️ Gateway no detectado:</span>
                      <span className="text-[10px]">{gatewayResponse.message}</span>
                    </div>
                  ) : (
                    <div className="flex flex-col gap-1">
                      <div className="flex justify-between items-center bg-blue-50 text-blue-800 px-2 py-1 rounded font-semibold text-[10px]">
                        <span>SERVIDOR PROCESADOR:</span>
                        <span className="bg-blue-600 text-white rounded px-1.5 py-0.2">{gatewayResponse.processed_by}</span>
                      </div>
                      <div className="text-[10px] text-gray-500 font-mono mt-1">
                        ID del Evento: <span className="text-gray-800 font-semibold">{gatewayResponse.event_id}</span>
                        <br />
                        Estado: <span className="text-green-600 font-semibold">{gatewayResponse.message}</span>
                      </div>
                    </div>
                  )
                ) : (
                  <span className="text-gray-400">Sin respuesta del Gateway.</span>
                )}
              </div>

              {/* Resumen del Pedido */}
              <div className="flex flex-col gap-1 border-t border-gray-100 pt-3">
                <span className="font-bold text-gray-800 block mb-1">Resumen del pedido:</span>
                <div className="flex justify-between text-gray-600">
                  <span>{selectedQty}x {productName}</span>
                  <span className="font-semibold text-gray-900">$ {(productPrice * selectedQty).toLocaleString('es-AR')}</span>
                </div>
                <div className="flex justify-between text-gray-600">
                  <span>Envío</span>
                  <span className="text-green-600 font-semibold">Gratis</span>
                </div>
                <div className="flex justify-between text-gray-900 font-bold border-t border-gray-200 pt-2 mt-1">
                  <span>Total abonado</span>
                  <span>$ {(productPrice * selectedQty).toLocaleString('es-AR')}</span>
                </div>
              </div>

            </div>

            {/* Footer éxito */}
            <div className="p-4 bg-gray-50 border-t border-gray-200 flex justify-end">
              <button
                onClick={() => {
                  setShowCheckout(false);
                  setLoadingStep(0);
                  setGatewayResponse(null);
                }}
                className="px-5 py-2 bg-gray-900 hover:bg-black text-white rounded font-bold text-xs shadow-sm cursor-pointer"
              >
                Volver a la publicación
              </button>
            </div>

          </div>
        </div>
      )}

    </div>
  );
}
