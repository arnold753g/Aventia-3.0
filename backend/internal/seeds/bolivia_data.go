package seeds

// SeedDepartamento estructura para organizar departamentos y provincias
type SeedDepartamento struct {
	Nombre     string
	Provincias []string
}

// Bolivia contiene todos los departamentos y provincias de Bolivia
var Bolivia = []SeedDepartamento{
	{
		Nombre: "Chuquisaca",
		Provincias: []string{
			"Oropeza",
			"Azurduy",
			"Zudáñez",
			"Tomina",
			"Hernando Siles",
			"Yamparáez",
			"Nor Cinti",
			"Sud Cinti",
			"Belisario Boeto",
			"Luis Calvo",
		},
	},
	{
		Nombre: "La Paz",
		Provincias: []string{
			"Abel Iturralde",
			"Aroma",
			"Bautista Saavedra",
			"Caranavi",
			"Eliodoro Camacho",
			"Franz Tamayo",
			"Gualberto Villarroel",
			"Ingavi",
			"Inquisivi",
			"José Manuel Pando",
			"Larecaja",
			"Loayza",
			"Los Andes",
			"Manco Kapac",
			"Muñecas",
			"Nor Yungas",
			"Omasuyos",
			"Pacajes",
			"Pedro Domingo Murillo",
			"Sud Yungas",
		},
	},
	{
		Nombre: "Cochabamba",
		Provincias: []string{
			"Arani",
			"Arque",
			"Ayopaya",
			"Bolívar",
			"Capinota",
			"Carrasco",
			"Cercado",
			"Chapare",
			"Esteban Arce",
			"Germán Jordán",
			"Mizque",
			"Narciso Campero",
			"Punata",
			"Quillacollo",
			"Tapacarí",
			"Tiraque",
		},
	},
	{
		Nombre: "Oruro",
		Provincias: []string{
			"Cercado",
			"Eduardo Avaroa",
			"Ladislao Cabrera",
			"Litoral",
			"Mejillones",
			"Nor Carangas",
			"Pantaleón Dalence",
			"Poopó",
			"Sabaya",
			"Sajama",
			"San Pedro de Totora",
			"Saucarí",
			"Sebastián Pagador",
			"Sud Carangas",
			"Tomás Barrón",
		},
	},
	{
		Nombre: "Potosí",
		Provincias: []string{
			"Alonso de Ibáñez",
			"Antonio Quijarro",
			"Bernardino Bilbao",
			"Charcas",
			"Chayanta",
			"Cornelio Saavedra",
			"Daniel Campos",
			"Enrique Baldivieso",
			"José María Linares",
			"Modesto Omiste",
			"Nor Chichas",
			"Nor Lípez",
			"Rafael Bustillo",
			"Sud Chichas",
			"Sud Lípez",
			"Tomás Frías",
		},
	},
	{
		Nombre: "Tarija",
		Provincias: []string{
			"Aniceto Arce",
			"Burnet O'Connor",
			"Cercado",
			"Eustaquio Méndez",
			"Gran Chaco",
			"José María Avilés",
		},
	},
	{
		Nombre: "Santa Cruz",
		Provincias: []string{
			"Andrés Ibáñez",
			"Ángel Sandóval",
			"Chiquitos",
			"Cordillera",
			"Florida",
			"Germán Busch",
			"Guarayos",
			"Ichilo",
			"José Miguel de Velasco",
			"Manuel María Caballero",
			"Ñuflo de Chávez",
			"Obispo Santistevan",
			"Sara",
			"Vallegrande",
			"Warnes",
		},
	},
	{
		Nombre: "Beni",
		Provincias: []string{
			"Cercado",
			"Antonio Vaca Díez",
			"General José Ballivián",
			"Iténez",
			"Mamoré",
			"Marbán",
			"Moxos",
			"Yacuma",
		},
	},
	{
		Nombre: "Pando",
		Provincias: []string{
			"Abuná",
			"Federico Román",
			"Madre de Dios",
			"Manuripi",
			"Nicolás Suárez",
		},
	},
}

// Totales para validación
const (
	TotalDepartamentos = 9
	TotalProvincias    = 111
)