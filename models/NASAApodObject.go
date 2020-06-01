package models

// NASAApodObject model
type NASAApodObject struct {
	Copyright       string `json:"copyright,omitempty"`
	Date            string `json:"date"`
	Explanation     string `json:"explanation"`
	Hdurl           string `json:"hdurl,omitempty"`
	Media_type      string `json:"media_type"`
	Service_version string `json:"service_version"`
	Title           string `json:"title"`
	Url             string `json:"url"`
}

/*
Example of an object returned by NASA
{
	copyright: "Göran Strand"
	date: "2020-05-31"
	explanation: "It was bright and green and stretched across the sky. This striking aurora display was captured in 2016 just outside of Östersund, Sweden. Six photographic fields were merged to create the featured panorama spanning almost 180 degrees.  Particularly striking aspects of this aurora include its sweeping arc-like shape and its stark definition.  Lake Storsjön is seen in the foreground, while several familiar constellations and the star Polaris are visible through the aurora, far in the background.  Coincidently, the aurora appears to avoid the Moon visible on the lower left.  The aurora appeared a day after a large hole opened in the Sun's corona allowing particularly energetic particles to flow out into the Solar System.  The green color of the aurora is caused by oxygen atoms recombining with ambient electrons high in the Earth's atmosphere."
	hdurl: "https://apod.nasa.gov/apod/image/2005/AuroraSweden_Strand_1500.jpg"
	media_type: "image"
	service_version: "v1"
	title: "Aurora over Sweden"
	url: "https://apod.nasa.gov/apod/image/2005/AuroraSweden_Strand_960.jpg"
}
*/
