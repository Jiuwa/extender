package controller

import (
	"log"
	"strings"
	
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (
	PodNameFitPred        = "PodNameFitPred"
	PodNameFitPredFailMsg = "Pod name don't fit "
)

var predicatesFuncs = map[string]FitPredicate{
	PodNameFitPred: PodNameFitPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{PodNameFitPred}

func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	pod := args.Pod

	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	failReasons := []string{}
	for _, predicateKey := range predicatesSorted {
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

func PodNameFitPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	var valid bool
	valid = (len(pod.Name) > 30)
	if valid {
		log.Printf("pod %v/%v length is %d, node length is over 10, fit", pod.Name, pod.Namespace, len(pod.Name))
		return true, nil, nil
	}
	log.Printf("pod %v/%v length is %d,  node length is not over 10, unfit\n", pod.Name, pod.Namespace, len(pod.Name))
	return false, []string{PodNameFitPredFailMsg}, nil
}
