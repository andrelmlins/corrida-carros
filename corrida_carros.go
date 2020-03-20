package main

import "fmt"
import "time"

var velocidade_maxima int = 3
var tamanho_pista int = 20
var count_faixas int = 3
var count_carros int = 3

var acelerar = make(chan int, 1)
var freiar = make(chan int, 1)
var manter = make(chan int, 1)
var mudarfaixa = make(chan MudarFaixa, 1)
var consultar = make(chan Consulta, 1)
var retornoConsultar = make(chan int, 1)
var ganhou = make(chan int, 1)
var terminou = make(chan int, 1)

var(
	obstaculo1 = Obstaculo{6,2}
	obstaculo2 = Obstaculo{9,1}
)

var vetorObstaculos = [2]Obstaculo{obstaculo1, obstaculo2}
	
type Carro struct {
	id int
	velocidade int
	posicao int
	faixa int
}

type Obstaculo struct {
	posicao int
	faixa int
}

type MudarFaixa struct {
	id int
	faixa int
}

type Consulta struct {
	faixa int
	posicao int
}

func pausa(){
	time.Sleep(100*time.Millisecond)	
}

func carros(id int, c chan int){
	for i:=0; i<10; i++ {
		c <- id
		pausa();
	}
}

func existeObstaculo(posicao int, posicaoPrevendo int, faixa int) bool {
	for i:=0; i<2; i++ {
		if(vetorObstaculos[i].posicao > posicao &&
		vetorObstaculos[i].posicao <= posicaoPrevendo &&
		vetorObstaculos[i].faixa == faixa){
			return true
		}
	}
	return false
}

func atualizaCarro(car *Carro, velocidade int, faixa int, posicao int){
	car.velocidade = velocidade
	car.posicao = posicao
	car.faixa = faixa
}

func consultarCarro(carros [][]int, faixa int, posicao int) (int){
	for i := 0; i < len(carros); i++ {
		if(carros[i][3]==faixa && carros[i][2]==posicao){
			return carros[i][0]
		}
	}

	return 0
}

func acelerarCarro(carros [][]int, id int) ([][]int){
	for i := 0; i < len(carros); i++ {
		if(carros[i][0]==id){
			velocidade:=carros[i][1]
			posicao:=carros[i][2]

			if(velocidade+1 <= velocidade_maxima && posicao+velocidade+1 <= tamanho_pista){
				velocidade++
				posicao+=velocidade
			} else if (posicao+velocidade+1 <= tamanho_pista){
				posicao+=velocidade
			}

			carros[i][1] = velocidade
			carros[i][2] = posicao

			break
		}
	}

	return carros
}

func freiarCarro(carros [][]int, id int) ([][]int){
	for i := 0; i < len(carros); i++ {
		if(carros[i][0]==id){
			velocidade:=carros[i][1]
			posicao:=carros[i][2]

			if(velocidade-1 >=0 && posicao+velocidade-1 <= tamanho_pista){
				velocidade--
				posicao+=velocidade
			} else if (posicao+velocidade+1 <= tamanho_pista){
				posicao+=velocidade
			}

			carros[i][1] = velocidade
			carros[i][2] = posicao

			break
		}
	}

	return carros
}

func manterCarro(carros [][]int, id int) ([][]int){
	for i := 0; i < len(carros); i++ {
		if(carros[i][0]==id){
			velocidade:=carros[i][1]
			posicao:=carros[i][2]

			if(posicao+velocidade <= tamanho_pista){
				posicao+=velocidade
			}

			carros[i][1] = velocidade
			carros[i][2] = posicao

			break
		}
	}

	return carros
}

func mudarfaixaCarro(carros [][]int, id int, faixa_nova int) ([][]int){
	for i := 0; i < len(carros); i++ {
		if(carros[i][0]==id){
			velocidade:=carros[i][1]
			posicao:=carros[i][2]
			faixa:=carros[i][3]

			if(velocidade-1>=0 && posicao+velocidade-1<=tamanho_pista){
				velocidade--
				posicao+=velocidade
				faixa = faixa_nova
			} else if(posicao+velocidade-1<=tamanho_pista) {
				posicao+=velocidade
			}

			carros[i][1] = velocidade
			carros[i][2] = posicao
			carros[i][3] = faixa

			break
		}
	}

	return carros
}

func removerCarro(carros [][]int, id int) ([][]int){
	for i := 0; i < len(carros); i++ {
		if(carros[i][0]==id){
			return append(carros[:i], carros[i+1:]...)
		}
	}

	return carros
}		

func mudarFaixaVelocidade(id int, velocidade int, faixa int, posicao int) (int, int, int){
	if(faixa < count_faixas && faixa > 1){
		consulta := new(Consulta)
		consulta.faixa = faixa+1
		consulta.posicao = posicao+velocidade-1
		consultar <- *consulta
		id_ladoD := <- retornoConsultar

		consulta1 := new(Consulta)
		consulta1.faixa = faixa-1
		consulta1.posicao = posicao+velocidade-1
		consultar <- *consulta1
		id_ladoE := <- retornoConsultar

		if(faixa+1 != count_faixas && id_ladoD == 0){
			return velocidade-1,faixa+1,posicao+velocidade-1
		} else if (id_ladoE == 0){
			send := new(MudarFaixa)
			send.id = id
			send.faixa = faixa-1
			mudarfaixa <- *send
			return velocidade-1, faixa-1, posicao+velocidade-1
		} else if(id_ladoD == 0){
			send := new(MudarFaixa)
			send.id = id
			send.faixa = faixa+1
			mudarfaixa <- *send
			return velocidade-1, faixa+1, posicao+velocidade-1
		}else{
			freiar <- id
			return velocidade-1, faixa, posicao+velocidade-1
		}
	} else if(faixa < count_faixas){
		consulta := new(Consulta)
		consulta.faixa = faixa+1
		consulta.posicao = posicao+velocidade-1
		consultar <- *consulta
		id_ladoD := <- retornoConsultar

		if(id_ladoD == 0){
			send := new(MudarFaixa)
			send.id = id
			send.faixa = faixa+1
			mudarfaixa <- *send
			return velocidade-1, faixa+1, posicao+velocidade-1
		} else{
			freiar <- id
			return velocidade-1, faixa, posicao+velocidade-1
		}
	} else if(faixa > 1){
		consulta := new(Consulta)
		consulta.faixa = faixa-1
		consulta.posicao = posicao+velocidade-1
		consultar <- *consulta
		id_ladoE := <- retornoConsultar

		if (id_ladoE == 0){
			send := new(MudarFaixa)
			send.id = id
			send.faixa = faixa-1
			mudarfaixa <- *send
			return velocidade-1, faixa-1, posicao+velocidade-1
		} else{
			freiar <- id
			return velocidade-1, faixa, posicao+velocidade-1
		}
	}else{
		freiar <- id
		return velocidade-1, faixa, posicao+velocidade-1
	}
}

func carro(id int,velocidade int,posicao int,faixa int){
	for {
		fmt.Println("Andar,", id, velocidade, posicao, faixa)
		if(posicao < tamanho_pista){
			if(existeObstaculo(posicao, posicao+velocidade+1, faixa) && velocidade > 0){
				velocidade,posicao,faixa = mudarFaixaVelocidade(id, velocidade, faixa, posicao)
			} else {
				consulta := new(Consulta)
				consulta.faixa = faixa
				consulta.posicao = posicao+velocidade+1
				consultar <- *consulta
				id_frente := <- retornoConsultar

				if(id_frente == 0){
					if(velocidade < velocidade_maxima){
						acelerar <- id
						velocidade+=1
						posicao+=velocidade
					}else{
						manter <- id
						posicao+=velocidade
					}
				} else if(velocidade > 0){
					velocidade,posicao,faixa = mudarFaixaVelocidade(id, velocidade, faixa, posicao)
				} else{
					acelerar <- id
					velocidade+=1
					posicao+=velocidade+1
				}
			}
		} else {
			ganhou  <- id
			return;
		}
	}
}

func monitor(carros [][]int){
	for {
		var id int
		select {
			case id = <- acelerar :
				carros = acelerarCarro(carros,id)
				fmt.Println("Acelerar,", id)
			case id = <- freiar :
				carros = freiarCarro(carros,id)
				fmt.Println("Freiar,", id)
			case id = <- manter :
				carros = manterCarro(carros,id)
				fmt.Println("Manter,", id)
			case idfaixa := <- mudarfaixa :
				carros = mudarfaixaCarro(carros,idfaixa.id,idfaixa.faixa)
				fmt.Println("Mudar Faixa,", idfaixa)
			case faixaposicao := <- consultar :
				retornoConsultar <- consultarCarro(carros,faixaposicao.faixa,faixaposicao.posicao)
			case id = <- ganhou :
				carros = removerCarro(carros,id)
				fmt.Println("Ganhou,", id)
				if(len(carros)==0) {
					terminou <- 1
				}
		}
	}
}

func main() {
	carros := make([][]int, count_carros)

	for i := 0; i < count_carros; i++ {
		carros[i] = make([]int, 4)
		carros[i][0] = i+1 //ID
		carros[i][1] = 0 //VELOCIDADE
		carros[i][2] = 1 //POSIÇÃO
		carros[i][3] = i+1 //FAIXA
	}

	go monitor(carros)
	
	for i := 0; i < count_carros; i++ {
		go carro(i+1,0,1,i+1)
	}

	<-terminou
	fmt.Println("Terminou")
}
