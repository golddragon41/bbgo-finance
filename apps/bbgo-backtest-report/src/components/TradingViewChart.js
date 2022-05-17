import React, {useEffect, useRef, useState} from 'react';
import {tsvParse} from "d3-dsv";

// https://github.com/tradingview/lightweight-charts/issues/543
// const createChart = dynamic(() => import('lightweight-charts'));
import {createChart} from 'lightweight-charts';

// const parseDate = timeParse("%Y-%m-%d");

const parseKline = () => {
  return (d) => {
    d.startTime = new Date(Number(d.startTime) * 1000);
    d.endTime = new Date(Number(d.endTime) * 1000);
    d.time = d.startTime.getTime() / 1000;

    for (const key in d) {
      // convert number fields
      if (Object.prototype.hasOwnProperty.call(d, key)) {
        switch (key) {
          case "open":
          case "high":
          case "low":
          case "close":
          case "volume":
            d[key] = +d[key];
            break
        }
      }
    }

    return d;
  };
};


const parseOrder = () => {
  return (d) => {
    for (const key in d) {
      // convert number fields
      if (Object.prototype.hasOwnProperty.call(d, key)) {
        switch (key) {
          case "order_id":
          case "price":
          case "quantity":
            d[key] = +d[key];
            break;
          case "time":
            d[key] = new Date(d[key]);
            break;
        }
      }
    }
    return d;
  };
}

const parsePosition = () => {
  return (d) => {
    for (const key in d) {
      // convert number fields
      if (Object.prototype.hasOwnProperty.call(d, key)) {
        switch (key) {
          case "accumulated_profit":
          case "average_cost":
          case "quote":
          case "base":
            d[key] = +d[key];
            break
          case "time":
            d[key] = new Date(d[key]);
            break
        }
      }
    }
    return d;
  };
}


const fetchPositionHistory = (setter) => {
  return fetch(
    `/data/bollmaker:ETHUSDT-position.tsv`,
  )
    .then((response) => response.text())
    .then((data) => tsvParse(data, parsePosition()))
    // .then((data) => tsvParse(data))
    .then((data) => {
      setter(data);
    })
    .catch((e) => {
      console.error("failed to fetch orders", e)
    });
};

const fetchOrders = (setter) => {
  return fetch(
    `/data/orders.tsv`,
  )
    .then((response) => response.text())
    .then((data) => tsvParse(data, parseOrder()))
    // .then((data) => tsvParse(data))
    .then((data) => {
      setter(data);
    })
    .catch((e) => {
      console.error("failed to fetch orders", e)
    });
}

const ordersToMarkets = (orders) => {
  const markers = [];
  // var markers = [{ time: data[data.length - 48].time, position: 'aboveBar', color: '#f68410', shape: 'circle', text: 'D' }];
  for (let i = 0; i < orders.length; i++) {
    let order = orders[i];
    switch (order.side) {
      case "BUY":
        markers.push({
          time: order.time.getTime() / 1000.0,
          position: 'belowBar',
          color: '#239D10',
          shape: 'arrowDown',
          // text: 'Buy @ ' + order.price
          text: 'B',
        });
        break;
      case "SELL":
        markers.push({
          time: order.time.getTime() / 1000.0,
          position: 'aboveBar',
          color: '#e91e63',
          shape: 'arrowDown',
          // text: 'Sell @ ' + order.price
          text: 'S',
        });
        break;
    }
  }
  return markers;
};

function fetchKLines(symbol, interval, setter) {
  return fetch(
    `/data/klines/${symbol}-${interval}.tsv`,
  )
    .then((response) => response.text())
    .then((data) => tsvParse(data, parseKline()))
    // .then((data) => tsvParse(data))
    .then((data) => {
      setter(data);
    })
    .catch((e) => {
      console.error("failed to fetch klines", e)
    });
}

const klinesToVolumeData = (klines) => {
  const volumes = [];

  for (let i = 0 ; i < klines.length ; i++) {
    const kline = klines[i];
    volumes.push({
      time: (kline.startTime.getTime() / 1000),
      value: kline.volume,
    })
  }

  return volumes;
}

const positionAverageCostHistoryToLineData = (hs) => {
  const avgCosts = [];
  for (let i = 0; i < hs.length; i++) {
    let pos = hs[i];

    if (i > 0 && pos.average_cost == hs[i-1].average_cost) {
      continue;
    }

    if (pos.base == 0) {
      avgCosts.push({
        time: pos.time.getTime() / 1000,
        value: 0,
      });
    } else {
      avgCosts.push({
        time: pos.time.getTime() / 1000,
        value: pos.average_cost,
      });
    }


  }
  return avgCosts;
}

const TradingViewChart = (props) => {
  const ref = useRef();
  const [data, setData] = useState(null);
  const [orders, setOrders] = useState(null);
  const [markers, setMarkers] = useState(null);
  const [positionHistory, setPositionHistory] = useState(null);

  useEffect(() => {
    if (!ref.current || ref.current.children.length > 0) {
      return;
    }

    if (!data || !orders || !markers || !positionHistory) {
      fetchKLines('ETHUSDT', '5m', setData).then(() => {
        fetchOrders((orders) => {
          setOrders(orders);

          const markers = ordersToMarkets(orders);
          setMarkers(markers);
        });
        fetchPositionHistory(setPositionHistory)
      })
      return;
    }

    console.log("createChart")
    const chart = createChart(ref.current, {
      width: 800,
      height: 200,
      timeScale: {
        timeVisible: true,
        borderColor: '#D1D4DC',
      },
      rightPriceScale: {
        borderColor: '#D1D4DC',
      },
      layout: {
        backgroundColor: '#ffffff',
        textColor: '#000',
      },
      grid: {
        horzLines: {
          color: '#F0F3FA',
        },
        vertLines: {
          color: '#F0F3FA',
        },
      },
    });

    const series = chart.addCandlestickSeries({
      upColor: 'rgb(38,166,154)',
      downColor: 'rgb(255,82,82)',
      wickUpColor: 'rgb(38,166,154)',
      wickDownColor: 'rgb(255,82,82)',
      borderVisible: false,
    });
    series.setData(data);
    series.setMarkers(markers);

    const lineSeries = chart.addLineSeries();
    const costLine = positionAverageCostHistoryToLineData(positionHistory);
    lineSeries.setData(costLine);

    const volumeData = klinesToVolumeData(data);
    const volumeSeries = chart.addHistogramSeries({
      color: '#182233',
      lineWidth: 2,
      priceFormat: {
        type: 'volume',
      },
      overlay: true,
      scaleMargins: {
        top: 0.8,
        bottom: 0,
      },
    });
    volumeSeries.setData(volumeData);

  }, [ref.current, data])
  return <div>
    <div ref={ref}>

    </div>
  </div>;
};

export default TradingViewChart;
