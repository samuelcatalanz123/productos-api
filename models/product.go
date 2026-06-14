package models

import "encoding/xml"

// Product representa un producto.
// Las "etiquetas" (`json`, `xml`, `yaml`) le dicen a Go cómo nombrar cada
// campo cuando lo convierta a cada formato.
type Product struct {
	XMLName xml.Name `json:"-" yaml:"-" xml:"producto"`
	ID      int      `json:"id" xml:"id" yaml:"id"`
	Nombre  string   `json:"nombre" xml:"nombre" yaml:"nombre"`
	Precio  float64  `json:"precio" xml:"precio" yaml:"precio"`
}

// ProductList envuelve una lista de productos (para que el XML tenga una raíz).
type ProductList struct {
	XMLName   xml.Name  `json:"-" yaml:"-" xml:"productos"`
	Productos []Product `json:"productos" xml:"producto" yaml:"productos"`
}
