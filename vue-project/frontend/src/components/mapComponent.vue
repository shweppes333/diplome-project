<template>
  <div class="map-container" id="mapContainer" style="height: 300px;">
    <div id="yandexMap" style="height:100%;"></div>
  </div>
</template>

<script>
export default {
  name: 'MapComponent',
  mounted() {
    this.initMap()
  },
  methods: {
    initMap() {
      const script = document.createElement('script')
      script.src = 'https://api-maps.yandex.ru/2.1/?apikey=ваш_ключ_api&lang=ru_RU'
      script.onload = () => {
        window.ymaps.ready(() => {
          const map = new window.ymaps.Map("yandexMap", {
            center: [55.751574, 37.633856],
            zoom: 14,
            controls: ['zoomControl', 'fullscreenControl']
          })
          
          const placemark = new window.ymaps.Placemark([55.751574, 37.633856], {
            balloonContent: "HandMadeStudio<br>ул. Казакова 8, мастерская 12"
          }, {
            preset: 'islands#redIcon'
          })
          
          map.geoObjects.add(placemark)
        })
      }
      document.head.appendChild(script)
    }
  }
}
</script>