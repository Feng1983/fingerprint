package main

import (
	"fmt"
	"github.com/xlvector/hector"
	"github.com/xlvector/hector/core"
	"github.com/xlvector/hector/svm"
	"flag"
	//"strconv"
)

func Prepareparams()(string,string,string,map[string]string){
	params := make(map[string]string)
	//flag.Parse()
	train_path:= flag.String("train","data.txt","path of train file")
	test_path := flag.String("test","test.txt","path of test")
	method	  := flag.String("method","knn","algorithm name")
	fmt.Println(*method)
	fmt.Println(*train_path)
	fmt.Println(*test_path)
	params["k"]="4"
	return *train_path,*test_path,*method,params
}

func SplitFile(dataSet *core.DataSet,total, part int)(*core.DataSet, *core.DataSet){
	train:= core.NewDataSet()
	test:= core.NewDataSet()

	for i, sample:= range dataSet.Samples{
		if i% total == part{
			test.AddSample(sample)
		}else{
			train.AddSample(sample)
		}
	}
	return train, test
}
func typecase(general interface{}){
	 switch general.(type){
        case *core.Sample:
            fmt.Println("sample")
		case *svm.KNN:
			fmt.Println("knn")
        default:
            fmt.Println("Unknowtype")
    }
}
func main(){
	train,test,method, params:= Prepareparams()
	classifier := hector.GetClassifier(method)
	fmt.Println(test)
	dataset:= core.NewDataSet()
	dataset.Load(train, -1)
	trainset,testset:= SplitFile(dataset, 10,2)
	fmt.Println(len(trainset.Samples))
	fmt.Println("test set is : ",len(testset.Samples))
	for _,v := range testset.Samples{
		fmt.Println("sample...",v)
		for _,f := range v.Features{
			fmt.Println(f.Id,"-> ",f.Value)
		}
	}
	
	classifier.Init(params)
	//auc, predict := hector.AlgorithmRunOnDataSet(classifier, trainset, testset, "", params)
	//fmt.Println(auc)
	//fmt.Println(predict)
	classifier.Train(trainset)
	if knn,ok :=classifier.(*svm.KNN);!ok{
        fmt.Println(knn.GetSv())
		fmt.Println("is not KNN*")
    }
	//fmt.Println(knn)
	auc, predict := hector.AlgorithmRunOnDataSet(classifier, trainset, testset, "", params)
	fmt.Println(predict,auc)
}
