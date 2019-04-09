package cophymaru

import (
	"fmt"
	"math"
	"os"
)

func postorder(curnode *Node) {
	for _, chld := range curnode.CHLD {
		postorder(chld)
	}
	fmt.Println(curnode.NAME, curnode.CONTRT)
}

//BMPruneRooted will prune BM branch lens and PICs down to a rooted node
//root node should be a real (ie. bifurcating) root
func BMPruneRooted(n *Node) {
	for _, chld := range n.CHLD {
		BMPruneRooted(chld)
	}
	n.PRNLEN = n.LEN
	nchld := len(n.CHLD)
	if nchld != 0 { //&& n.MRK == false {
		var tempChar float64
		if nchld != 2 {
			fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
		}
		c0 := n.CHLD[0]
		c1 := n.CHLD[1]
		bot := ((1.0 / c0.PRNLEN) + (1.0 / c1.PRNLEN))
		n.PRNLEN += 1.0 / bot
		for i := range n.CHLD[0].CONTRT {
			tempChar = (((1 / c0.PRNLEN) * c1.CONTRT[i]) + ((1 / c1.PRNLEN) * c0.CONTRT[i])) / bot
			n.CONTRT[i] = tempChar
		}
	}
}

//AncBMPruneRooted will prune BM branch lens and PICs down to a rooted node
//root node should be a real (ie. bifurcating) root
func AncBMPruneRooted(n *Node) {
	for _, chld := range n.CHLD {
		AncBMPruneRooted(chld)
	}
	n.LSPRNLEN = n.LSLEN
	nchld := len(n.CHLD)
	if nchld != 0 { //&& n.MRK == false {
		var tempChar float64
		if nchld != 2 {
			fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
		}
		c0 := n.CHLD[0]
		c1 := n.CHLD[1]
		bot := ((1.0 / c0.LSPRNLEN) + (1.0 / c1.LSPRNLEN))
		n.LSPRNLEN += 1.0 / bot
		for i := range n.CHLD[0].CONTRT {
			tempChar = (((1 / c0.LSPRNLEN) * c1.CONTRT[i]) + ((1 / c1.LSPRNLEN) * c0.CONTRT[i])) / bot
			n.CONTRT[i] = tempChar
		}
	}
}

//AncBMPruneRootedSingle will prune BM branch lens and calculate PIC of a single trait down to a rooted node
//root node should be a real (ie. bifurcating) root
func AncBMPruneRootedSingle(n *Node, i int) {
	for _, chld := range n.CHLD {
		AncBMPruneRootedSingle(chld, i)
	}
	n.LSPRNLEN = n.LSLEN
	nchld := len(n.CHLD)
	if nchld != 0 { //&& n.MRK == false {
		if nchld != 2 {
			fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
		}
		c0 := n.CHLD[0]
		c1 := n.CHLD[1]
		bot := ((1.0 / c0.LSPRNLEN) + (1.0 / c1.LSPRNLEN))
		n.LSPRNLEN += 1.0 / bot
		tempCharacter := (((1 / c0.LSPRNLEN) * c1.CONTRT[i]) + ((1 / c1.LSPRNLEN) * c0.CONTRT[i])) / bot
		n.CONTRT[i] = tempCharacter
	}
}

//BMPruneRootedSingle will prune BM branch lens and calculate PIC of a single trait down to a rooted node
//root node should be a real (ie. bifurcating) root
func BMPruneRootedSingle(n *Node, i int) {
	for _, chld := range n.CHLD {
		BMPruneRootedSingle(chld, i)
	}
	n.PRNLEN = n.LEN
	nchld := len(n.CHLD)
	if nchld != 0 { //&& n.MRK == false {
		if nchld != 2 {
			fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
		}
		c0 := n.CHLD[0]
		c1 := n.CHLD[1]
		bot := ((1.0 / c0.PRNLEN) + (1.0 / c1.PRNLEN))
		n.PRNLEN += 1.0 / bot
		tempCharacter := (((1 / c0.PRNLEN) * c1.CONTRT[i]) + ((1 / c1.PRNLEN) * c0.CONTRT[i])) / bot
		n.CONTRT[i] = tempCharacter
	}
}

func debugParChld(tree *Node) {
	for _, ch := range tree.CHLD {
		fmt.Println(ch.NAME, "\t")
	}
}

//AssertUnrootedTree is a quick check to make sure the tree passed is unrooted
func AssertUnrootedTree(tree *Node) {
	if len(tree.CHLD) != 3 {
		fmt.Print("BRANCH LENGTHS MUST BE ITERATED ON AN UNROOTED TREE. THIS TREE IS ROOTED.")
		os.Exit(0)
	}
}

//IterateBMLengths will iteratively calculate the ML branch lengths for a particular topology, assuming that traits are fully sampled at the tips. should use MissingTraitsEM if there are missing sites.
func IterateBMLengths(tree *Node, niter int) {
	AssertUnrootedTree(tree)
	itercnt := 0
	for {
		calcBMLengths(tree)
		itercnt++
		if itercnt == niter {
			break
		}
	}
}

//AncMissingTraitsEM will iteratively calculate the ML branch lengths for a particular topology
func AncMissingTraitsEM(tree *Node, niter int) (nparam float64) {
	AssertUnrootedTree(tree)
	nodes := tree.PreorderArray()
	InitMissingValues(nodes)
	itercnt := 0
	for {
		AncCalcExpectedTraits(tree) //calculate Expected trait values
		ancCalcBMLengths(tree)      //maximize likelihood of branch lengths
		itercnt++
		if itercnt == niter {
			break
		}
	}
	nparam = 0.0
	for _, n := range nodes {
		if n.ANC == true && n.ISTIP == true {
			continue
		}
		nparam += 1.0
	}
	return
}

//MissingTraitsEM will iteratively calculate the ML branch lengths for a particular topology
func MissingTraitsEM(tree *Node, niter int) {
	AssertUnrootedTree(tree)
	nodes := tree.PreorderArray()
	InitMissingValues(nodes)
	itercnt := 0
	for {
		CalcExpectedTraits(tree) //calculate Expected trait values
		calcBMLengths(tree)      //maximize likelihood of branch lengths
		itercnt++
		if itercnt == niter {
			break
		}
	}
}

//ancCalcBMLengths will perform a single pass of the branch length ML estimation
func ancCalcBMLengths(tree *Node) {
	rnodes := tree.PreorderArray()
	lnode := 0
	for ind, newroot := range rnodes {
		if len(newroot.CHLD) == 0 {
			continue
		} else if newroot != rnodes[0] {
			tree = newroot.RerootLS(rnodes[lnode])
			lnode = ind
		}
		for _, cn := range tree.CHLD {
			AncBMPruneRooted(cn)
		}
		AncTritomyML(tree)
	}
	tree = rnodes[0].RerootLS(tree)
	//fmt.Println(tree.Newick(true))
}

//calcBMLengths will perform a single pass of the branch length ML estimation
func calcBMLengths(tree *Node) {
	rnodes := tree.PreorderArray()
	lnode := 0
	for ind, newroot := range rnodes {
		if len(newroot.CHLD) == 0 {
			continue
		} else if newroot != rnodes[0] {
			tree = newroot.Reroot(rnodes[lnode])
			lnode = ind
		}
		for _, cn := range tree.CHLD {
			AncBMPruneRooted(cn)
		}
		TritomyML(tree)
	}
	tree = rnodes[0].Reroot(tree)
	//fmt.Println(tree.Newick(true))
}

//PruneToStar will prune brlens and traits to a root
func PruneToStar(tree *Node) {
	for _, cn := range tree.CHLD {
		BMPruneRooted(cn)
	}
}

//CalcUnrootedLogLike will calculate the log-likelihood of an unrooted tree, while assuming that no sites have missing data.
func CalcUnrootedLogLike(tree *Node, startFresh bool) (chll float64) {
	chll = 0.0
	for _, ch := range tree.CHLD {
		curlike := 0.0
		CalcRootedLogLike(ch, &curlike, startFresh)
		chll += curlike
	}
	sitelikes := 0.0
	var tmpll float64
	var contrast, curVar float64
	for i := range tree.CHLD[0].CONTRT {
		tmpll = 0.
		if tree.CHLD[0].MIS[i] == false && tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == false { //do the standard calculation when no subtrees have missing traits
			contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[1].CONTRT[i]
			curVar = tree.CHLD[0].PRNLEN + tree.CHLD[1].PRNLEN
			tmpll = ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
			tmpPRNLEN := ((tree.CHLD[0].PRNLEN * tree.CHLD[1].PRNLEN) / (tree.CHLD[0].PRNLEN + tree.CHLD[1].PRNLEN))
			tmpChar := ((tree.CHLD[0].PRNLEN * tree.CHLD[1].CONTRT[i]) + (tree.CHLD[1].PRNLEN * tree.CHLD[0].CONTRT[i])) / curVar
			contrast = tmpChar - tree.CHLD[2].CONTRT[i]
			curVar = tree.CHLD[2].PRNLEN + tmpPRNLEN
			tmpll += ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
		}
		sitelikes += tmpll
	}
	chll += sitelikes
	return
}

//CalcRootedLogLike will return the BM likelihood of a tree assuming that no data are missing from the tips.
func CalcRootedLogLike(n *Node, nlikes *float64, startFresh bool) {
	for _, chld := range n.CHLD {
		CalcRootedLogLike(chld, nlikes, startFresh)
	}
	nchld := len(n.CHLD)
	if n.MRK == true && startFresh == false {
		if nchld != 0 {
			for _, l := range n.LL {
				*nlikes += l
			}
		}
	} else if n.MRK == false || startFresh == true {
		n.PRNLEN = n.LEN
		if nchld != 0 {
			if nchld != 2 {
				fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
			}
			c0 := n.CHLD[0]
			c1 := n.CHLD[1]
			curlike := float64(0.0)
			var tempChar float64
			for i := range n.CHLD[0].CONTRT {
				curVar := c0.PRNLEN + c1.PRNLEN
				contrast := c0.CONTRT[i] - c1.CONTRT[i]
				curlike += ((-0.5) * ((math.Log(2 * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2) / (curVar))))
				n.LL[i] = curlike
				tempChar = ((c0.PRNLEN * c1.CONTRT[i]) + (c1.PRNLEN * c0.CONTRT[i])) / (curVar)
				n.CONTRT[i] = tempChar
			}
			*nlikes += curlike
			tempBranchLength := n.LEN + ((c0.PRNLEN * c1.PRNLEN) / (c0.PRNLEN + c1.PRNLEN)) // need to calculate the prune length by adding the averaged lengths of the daughter nodes to the length
			n.PRNLEN = tempBranchLength                                                     // need to calculate the "prune length" by adding the length to the uncertainty

			n.MRK = true
		}
	}
}

//SitewiseLogLike will calculate the log-likelihood of an unrooted tree, while assuming that some sites have missing data. This can be used to calculate the likelihoods of trees that have complete trait sampling, but it will be slower than CalcRootedLogLike.
func SitewiseLogLike(tree *Node) (sitelikes []float64) {
	var tmpll float64
	ch1 := tree.CHLD[0]                     //.PostorderArray()
	ch2 := tree.CHLD[1]                     //.PostorderArray()
	ch3 := tree.CHLD[2]                     //.PostorderArray()
	for site := range tree.CHLD[0].CONTRT { //calculate log likelihood at each site
		tmpll = 0.
		calcRootedSiteLL(ch1, &tmpll, true, site)
		calcRootedSiteLL(ch2, &tmpll, true, site)
		calcRootedSiteLL(ch3, &tmpll, true, site)
		tmpll += calcUnrootedSiteLL(tree, site)
		//fmt.Println(site, tmpll)
		sitelikes = append(sitelikes, tmpll)
	}
	return
}

/*/WeightedUnrootedLogLike will calculate the log-likelihood of an unrooted tree, while assuming that no sites have missing data.
func WeightedUnrootedLogLike(tree *Node, startFresh bool, weights []float64) (chll float64) {
	chll = 0.0
	for _, ch := range tree.CHLD {
		n := ch.PostorderArray()
		for i := range ch.CONTRT {
			curlike := calcRootedSiteLL(n, startFresh, i)
			chll += curlike
		}
	}
	sitelikes := 0.0
	var tmpll float64
	var contrast, curVar float64
	for i := range tree.CHLD[0].CONTRT {
		tmpll = 0.
		if tree.CHLD[0].MIS[i] == false && tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == false { //do the standard calculation when no subtrees have missing traits
			contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[1].CONTRT[i]
			curVar = tree.CHLD[0].PRNLEN + tree.CHLD[1].PRNLEN
			tmpll = ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
			tmpPRNLEN := ((tree.CHLD[0].PRNLEN * tree.CHLD[1].PRNLEN) / (tree.CHLD[0].PRNLEN + tree.CHLD[1].PRNLEN))
			tmpChar := ((tree.CHLD[0].PRNLEN * tree.CHLD[1].CONTRT[i]) + (tree.CHLD[1].PRNLEN * tree.CHLD[0].CONTRT[i])) / curVar
			contrast = tmpChar - tree.CHLD[2].CONTRT[i]
			curVar = tree.CHLD[2].PRNLEN + tmpPRNLEN
			tmpll += ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
		}
		sitelikes += tmpll * weights[i]
	}
	chll += sitelikes
	return
}
*/

//MarkAll will mark all of the nodes in a tree ARRAY
func MarkAll(nodes []*Node) {
	for _, n := range nodes {
		n.MRK = true
	}
}

//calcUnrootedNodeLikes  will calculate the likelihood of an unrooted tree at each site (i) of the continuous character alignment
func calcUnrootedSiteLLParallel(tree *Node, i int) (tmpll float64) {
	var contrast, curVar float64
	log2pi := 1.8378770664093453
	if tree.CHLD[0].MIS[i] == false && tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == false { //do the standard calculation when no subtrees have missing traits
		contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[1].CONTRT[i]
		curVar = tree.CHLD[0].CONPRNLEN[i] + tree.CHLD[1].CONPRNLEN[i]
		tmpll = ((-0.5) * ((log2pi) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
		tmpCONPRNLEN := ((tree.CHLD[0].CONPRNLEN[i] * tree.CHLD[1].CONPRNLEN[i]) / (tree.CHLD[0].CONPRNLEN[i] + tree.CHLD[1].CONPRNLEN[i]))
		tmpChar := ((tree.CHLD[0].CONPRNLEN[i] * tree.CHLD[1].CONTRT[i]) + (tree.CHLD[1].CONPRNLEN[i] * tree.CHLD[0].CONTRT[i])) / curVar
		contrast = tmpChar - tree.CHLD[2].CONTRT[i]
		curVar = tree.CHLD[2].CONPRNLEN[i] + tmpCONPRNLEN
		tmpll += ((-0.5) * ((log2pi) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	} else if tree.CHLD[0].MIS[i] == false && tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == true { // do standard "rooted" calculation on CHLD[0] and CHLD [1] if CHLD[2] is missing
		contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[1].CONTRT[i]
		curVar = tree.CHLD[0].CONPRNLEN[i] + tree.CHLD[1].CONPRNLEN[i]
		tmpll = ((-0.5) * ((log2pi) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	} else if tree.CHLD[0].MIS[i] == false && tree.CHLD[2].MIS[i] == false && tree.CHLD[1].MIS[i] == true { // do standard "rooted" calculation on CHLD[0] and CHLD [2] if CHLD[1] is missing
		contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[2].CONTRT[i]
		curVar = tree.CHLD[0].CONPRNLEN[i] + tree.CHLD[2].CONPRNLEN[i]
		tmpll = ((-0.5) * ((log2pi) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	} else if tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == false && tree.CHLD[0].MIS[i] == true { // do standard "rooted" calculation on CHLD[1] and CHLD [2] if CHLD[0] is missing
		contrast = tree.CHLD[1].CONTRT[i] - tree.CHLD[2].CONTRT[i]
		curVar = tree.CHLD[1].CONPRNLEN[i] + tree.CHLD[2].CONPRNLEN[i]
		tmpll = ((-0.5) * ((log2pi) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	}
	return
}

//calcRootedSiteLL will return the BM likelihood of a tree assuming that no data are missing from the tips.
func calcRootedSiteLLParallel(n *Node, nlikes *float64, startFresh bool, site int) {
	for _, chld := range n.CHLD {
		calcRootedSiteLLParallel(chld, nlikes, startFresh, site)
	}
	nchld := len(n.CHLD)
	if n.MRK == true {
		if startFresh == false {
			if nchld != 0 {
				*nlikes += n.LL[site]
			}
		}
	}
	if n.MRK == false || startFresh == true {
		n.CONPRNLEN[site] = n.LEN
		log2pi := 1.8378770664093453
		if nchld != 0 {
			if nchld != 2 {
				fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
				os.Exit(0)
			}
			c0 := n.CHLD[0]
			c1 := n.CHLD[1]
			if c0.MIS[site] == false && c1.MIS[site] == false {
				curlike := float64(0.0)
				var tempChar float64
				curVar := (c0.CONPRNLEN[site]) + (c1.CONPRNLEN[site])
				contrast := c0.CONTRT[site] - c1.CONTRT[site]
				curlike += ((-0.5) * ((log2pi) + (math.Log(curVar)) + (math.Pow(contrast, 2) / (curVar))))
				tempChar = (((c0.CONPRNLEN[site]) * c1.CONTRT[site]) + ((c1.CONPRNLEN[site]) * c0.CONTRT[site])) / (curVar)
				n.CONTRT[site] = tempChar
				*nlikes += curlike
				tempBranchLength := n.CONPRNLEN[site] + (((c0.CONPRNLEN[site]) * (c1.CONPRNLEN[site])) / ((c0.CONPRNLEN[site]) + (c1.CONPRNLEN[site]))) // need to calculate the prune length by adding the averaged lengths of the daughter nodes to the length
				n.CONPRNLEN[site] = tempBranchLength                                                                                                    // need to calculate the "prune length" by adding the length to the uncertainty
				n.LL[site] = curlike
				//n.MRK = true
			} else if c0.MIS[site] == true && c1.MIS[site] == false {
				n.CONPRNLEN[site] += (c1.CONPRNLEN[site])
				n.CONTRT[site] = c1.CONTRT[site]
				n.LL[site] = 0.0
			} else if c1.MIS[site] == true && c0.MIS[site] == false {
				n.CONPRNLEN[site] += c0.CONPRNLEN[site]
				n.CONTRT[site] = c0.CONTRT[site]
				n.LL[site] = 0.0
			}
		}
	}
}

/*/calcRootedSiteLL will return the BM likelihood of a tree assuming that no data are missing from the tips.
func calcRootedSiteLLParallel(n *Node, nlikes *float64, startFresh bool, site int) {
	for _, chld := range n.CHLD {
		calcRootedSiteLLParallel(chld, nlikes, startFresh, site)
	}
	nchld := len(n.CHLD)
	if n.MRK == true {
		if startFresh == false {
			if nchld != 0 {
				*nlikes += n.LL[site]
			}
		}
	}
	if n.MRK == false || startFresh == true {
		n.CONPRNLEN[site] = n.LEN
		log2pi := 1.8378770664093453
		if nchld != 0 {
			if nchld != 2 {
				fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
				os.Exit(0)
			}
			c0 := n.CHLD[0]
			c1 := n.CHLD[1]
			curlike := float64(0.0)
			var tempChar float64
			curVar := c0.CONPRNLEN[site] + c1.CONPRNLEN[site]
			contrast := c0.CONTRT[site] - c1.CONTRT[site]
			curlike += ((-0.5) * ((log2pi) + (math.Log(curVar)) + (math.Pow(contrast, 2) / (curVar))))
			tempChar = ((c0.CONPRNLEN[site] * c1.CONTRT[site]) + (c1.CONPRNLEN[site] * c0.CONTRT[site])) / (curVar)
			n.CONTRT[site] = tempChar
			*nlikes += curlike
			tempBranchLength := n.CONPRNLEN[site] + ((c0.CONPRNLEN[site] * c1.CONPRNLEN[site]) / (c0.CONPRNLEN[site] + c1.CONPRNLEN[site])) // need to calculate the prune length by adding the averaged lengths of the daughter nodes to the length
			n.CONPRNLEN[site] = tempBranchLength                                                                                            // need to calculate the "prune length" by adding the length to the uncertainty
			n.LL[site] = curlike
			//n.MRK = true
		}
	}
}*/

func siteTreeLikeParallel(tree, ch1, ch2, ch3 *Node, startFresh bool, weights []float64, jobs <-chan int, results chan<- float64) {
	for site := range jobs {
		tmpll := 0.
		calcRootedSiteLLParallel(ch1, &tmpll, startFresh, site)
		calcRootedSiteLLParallel(ch2, &tmpll, startFresh, site)
		calcRootedSiteLLParallel(ch3, &tmpll, startFresh, site)
		tmpll += calcUnrootedSiteLLParallel(tree, site)
		tmpll = tmpll * weights[site]
		results <- tmpll
	}
}

//WeightedUnrootedLogLikeParallel will calculate the log-likelihood of an unrooted tree, while assuming that some sites have missing data. This can be used to calculate the likelihoods of trees that have complete trait sampling, but it will be slower than CalcRootedLogLike.
func WeightedUnrootedLogLikeParallel(tree *Node, startFresh bool, weights []float64, workers int) (sitelikes float64) {
	nsites := len(tree.CHLD[0].CONTRT)
	ch1 := tree.CHLD[0] //.PostorderArray()
	ch2 := tree.CHLD[1] //.PostorderArray()
	ch3 := tree.CHLD[2] //.PostorderArray()
	jobs := make(chan int, nsites)
	results := make(chan float64, nsites)
	for w := 0; w < workers; w++ {
		go siteTreeLikeParallel(tree, ch1, ch2, ch3, startFresh, weights, jobs, results)
	}

	for site := 0; site < nsites; site++ {
		jobs <- site
	}
	close(jobs)

	for site := 0; site < nsites; site++ {
		sitelikes += <-results
	}
	return
}

//WeightedUnrootedLogLike will calculate the log-likelihood of an unrooted tree, while assuming that some sites have missing data. This can be used to calculate the likelihoods of trees that have complete trait sampling, but it will be slower than CalcRootedLogLike.
func WeightedUnrootedLogLike(tree *Node, startFresh bool, weights []float64) (sitelikes float64) {
	sitelikes = 0.0
	var tmpll float64
	ch1 := tree.CHLD[0]                     //.PostorderArray()
	ch2 := tree.CHLD[1]                     //.PostorderArray()
	ch3 := tree.CHLD[2]                     //.PostorderArray()
	for site := range tree.CHLD[0].CONTRT { //calculate log likelihood at each site
		tmpll = 0.
		calcRootedSiteLL(ch1, &tmpll, startFresh, site)
		calcRootedSiteLL(ch2, &tmpll, startFresh, site)
		calcRootedSiteLL(ch3, &tmpll, startFresh, site)
		tmpll += calcUnrootedSiteLL(tree, site)
		tmpll = tmpll * weights[site]
		//fmt.Println(tmpll)
		sitelikes += tmpll
	}
	return
}

//MissingUnrootedLogLike will calculate the log-likelihood of an unrooted tree. It is deprecated and is only a convienience for testing LL calculation.
func MissingUnrootedLogLike(tree *Node, startFresh bool) (sitelikes float64) {
	sitelikes = 0.0
	var tmpll float64
	ch1 := tree.CHLD[0]                     //.PostorderArray()
	ch2 := tree.CHLD[1]                     //.PostorderArray()
	ch3 := tree.CHLD[2]                     //.PostorderArray()
	for site := range tree.CHLD[0].CONTRT { //calculate log likelihood at each site
		tmpll = 0.
		calcRootedSiteLL(ch1, &tmpll, startFresh, site)
		calcRootedSiteLL(ch2, &tmpll, startFresh, site)
		calcRootedSiteLL(ch3, &tmpll, startFresh, site)
		tmpll += calcUnrootedSiteLL(tree, site)
		sitelikes += tmpll
	}
	return
}

//calcUnrootedNodeLikes  will calculate the likelihood of an unrooted tree at each site (i) of the continuous character alignment
func calcUnrootedSiteLL(tree *Node, i int) (tmpll float64) {
	var contrast, curVar float64
	if tree.CHLD[0].MIS[i] == false && tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == false { //do the standard calculation when no subtrees have missing traits
		contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[1].CONTRT[i]
		curVar = tree.CHLD[0].PRNLEN + tree.CHLD[1].PRNLEN
		tmpll = ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
		tmpPRNLEN := ((tree.CHLD[0].PRNLEN * tree.CHLD[1].PRNLEN) / (tree.CHLD[0].PRNLEN + tree.CHLD[1].PRNLEN))
		tmpChar := ((tree.CHLD[0].PRNLEN * tree.CHLD[1].CONTRT[i]) + (tree.CHLD[1].PRNLEN * tree.CHLD[0].CONTRT[i])) / curVar
		contrast = tmpChar - tree.CHLD[2].CONTRT[i]
		curVar = tree.CHLD[2].PRNLEN + tmpPRNLEN
		tmpll += ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	} else if tree.CHLD[0].MIS[i] == false && tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == true { // do standard "rooted" calculation on CHLD[0] and CHLD [1] if CHLD[2] is missing
		contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[1].CONTRT[i]
		curVar = tree.CHLD[0].PRNLEN + tree.CHLD[1].PRNLEN
		tmpll = ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	} else if tree.CHLD[0].MIS[i] == false && tree.CHLD[2].MIS[i] == false && tree.CHLD[1].MIS[i] == true { // do standard "rooted" calculation on CHLD[0] and CHLD [2] if CHLD[1] is missing
		contrast = tree.CHLD[0].CONTRT[i] - tree.CHLD[2].CONTRT[i]
		curVar = tree.CHLD[0].PRNLEN + tree.CHLD[2].PRNLEN
		tmpll = ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	} else if tree.CHLD[1].MIS[i] == false && tree.CHLD[2].MIS[i] == false && tree.CHLD[0].MIS[i] == true { // do standard "rooted" calculation on CHLD[1] and CHLD [2] if CHLD[0] is missing
		contrast = tree.CHLD[1].CONTRT[i] - tree.CHLD[2].CONTRT[i]
		curVar = tree.CHLD[1].PRNLEN + tree.CHLD[2].PRNLEN
		tmpll = ((-0.5) * ((math.Log(2. * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2.) / (curVar))))
	}
	return
}

//MissingRootedLogLike will return the BM log-likelihood of a tree for a single site, pruning tips that have missing data
func MissingRootedLogLike(n *Node, startFresh bool) (sitelikes float64) {
	//nodes := n.PostorderArray()
	sitelikes = 0.
	for site := range n.CONTRT {
		calcRootedSiteLL(n, &sitelikes, startFresh, site)
	}
	return
}

//calcRootedSiteLL will return the BM likelihood of a tree assuming that no data are missing from the tips.
func calcRootedSiteLL(n *Node, nlikes *float64, startFresh bool, site int) {
	for _, chld := range n.CHLD {
		calcRootedSiteLL(chld, nlikes, startFresh, site)
	}
	nchld := len(n.CHLD)
	if n.MRK == true {
		if startFresh == false {
			if nchld != 0 {
				*nlikes += n.LL[site]
			}
		}
	}
	if n.MRK == false || startFresh == true {
		n.PRNLEN = n.LEN
		if nchld != 0 {
			if nchld != 2 {
				fmt.Println("This BM pruning algorithm should only be perfomed on fully bifurcating trees/subtrees! Check for multifurcations and singletons.")
				os.Exit(0)
			}
			c0 := n.CHLD[0]
			c1 := n.CHLD[1]
			curlike := float64(0.0)
			var tempChar float64
			curVar := c0.PRNLEN + c1.PRNLEN
			contrast := c0.CONTRT[site] - c1.CONTRT[site]
			curlike += ((-0.5) * ((math.Log(2 * math.Pi)) + (math.Log(curVar)) + (math.Pow(contrast, 2) / (curVar))))
			tempChar = ((c0.PRNLEN * c1.CONTRT[site]) + (c1.PRNLEN * c0.CONTRT[site])) / (curVar)
			n.CONTRT[site] = tempChar
			*nlikes += curlike
			tempBranchLength := n.LEN + ((c0.PRNLEN * c1.PRNLEN) / (c0.PRNLEN + c1.PRNLEN)) // need to calculate the prune length by adding the averaged lengths of the daughter nodes to the length
			n.PRNLEN = tempBranchLength                                                     // need to calculate the "prune length" by adding the length to the uncertainty
			n.LL[site] = curlike
			//n.MRK = true
		}
	}
}

//TritomyML will calculate the MLEs for the branch lengths of a tifurcating 3-taxon tree
func TritomyML(tree *Node) {
	ntraits := len(tree.CHLD[0].CONTRT)
	fntraits := float64(ntraits)
	var x1, x2, x3 float64
	sumV1 := 0.0
	sumV2 := 0.0
	sumV3 := 0.0
	for i := range tree.CHLD[0].CONTRT {
		x1 = tree.CHLD[0].CONTRT[i]
		x2 = tree.CHLD[1].CONTRT[i]
		x3 = tree.CHLD[2].CONTRT[i]
		sumV1 += ((x1 - x2) * (x1 - x3))
		sumV2 += ((x2 - x1) * (x2 - x3))
		sumV3 += ((x3 - x1) * (x3 - x2))
	}
	if sumV1 < 0.0 {
		sumV1 = 0.000001
		sumV2 = 0.0
		sumV3 = 0.0
		for i := range tree.CHLD[0].CONTRT {
			x1 = tree.CHLD[0].CONTRT[i]
			x2 = tree.CHLD[1].CONTRT[i]
			x3 = tree.CHLD[2].CONTRT[i]
			sumV2 += (x1 - x2) * (x1 - x2)
			sumV3 += (x1 - x3) * (x1 - x3)
		}
	} else if sumV2 < 0.0 {
		sumV1 = 0.0
		sumV2 = 0.00001
		sumV3 = 0.0
		for i := range tree.CHLD[0].CONTRT {
			x1 = tree.CHLD[0].CONTRT[i]
			x2 = tree.CHLD[1].CONTRT[i]
			x3 = tree.CHLD[2].CONTRT[i]
			sumV1 += (x2 - x1) * (x2 - x1)
			sumV3 += (x2 - x3) * (x2 - x3)
		}
	} else if sumV3 < 0.0 {
		sumV1 = 0.0
		sumV2 = 0.0
		sumV3 = 0.0001
		for i := range tree.CHLD[0].CONTRT {
			x1 = tree.CHLD[0].CONTRT[i]
			x2 = tree.CHLD[1].CONTRT[i]
			x3 = tree.CHLD[2].CONTRT[i]
			sumV1 += (x3 - x1) * (x3 - x1)
			sumV2 += (x3 - x2) * (x3 - x2)
		}
	}
	sumV1 = sumV1 / fntraits
	sumV2 = sumV2 / fntraits
	sumV3 = sumV3 / fntraits
	sumV1 = sumV1 - (tree.CHLD[0].PRNLEN - tree.CHLD[0].LEN)
	sumV2 = sumV2 - (tree.CHLD[1].PRNLEN - tree.CHLD[1].LEN)
	sumV3 = sumV3 - (tree.CHLD[2].PRNLEN - tree.CHLD[2].LEN)
	if sumV1 < 0. {
		sumV1 = 0.0001
	}
	if sumV2 < 0. {
		sumV2 = 0.0001
	}
	if sumV3 < 0. {
		sumV3 = 0.0001
	}
	tree.CHLD[0].LEN = sumV1
	tree.CHLD[1].LEN = sumV2
	tree.CHLD[2].LEN = sumV3
}

//AncTritomyML will calculate the MLEs for the branch lengths of a tifurcating 3-taxon tree assuming that direct ancestors may be in the tree
func AncTritomyML(tree *Node) {
	ntraits := len(tree.CHLD[0].CONTRT)
	fntraits := float64(ntraits)
	var x1, x2, x3 float64
	sumV1 := 0.0
	sumV2 := 0.0
	sumV3 := 0.0
	for i := range tree.CHLD[0].CONTRT {
		x1 = tree.CHLD[0].CONTRT[i]
		x2 = tree.CHLD[1].CONTRT[i]
		x3 = tree.CHLD[2].CONTRT[i]
		sumV1 += ((x1 - x2) * (x1 - x3))
		sumV2 += ((x2 - x1) * (x2 - x3))
		sumV3 += ((x3 - x1) * (x3 - x2))
	}
	if sumV1 < 0.0 || tree.CHLD[0].ANC == true && tree.CHLD[0].ISTIP == true {
		sumV1 = 0.0000000000001
		sumV2 = 0.0
		sumV3 = 0.0
		for i := range tree.CHLD[0].CONTRT {
			x1 = tree.CHLD[0].CONTRT[i]
			x2 = tree.CHLD[1].CONTRT[i]
			x3 = tree.CHLD[2].CONTRT[i]
			sumV2 += (x1 - x2) * (x1 - x2)
			sumV3 += (x1 - x3) * (x1 - x3)
		}
	} else if sumV2 < 0.0 || tree.CHLD[1].ANC == true && tree.CHLD[1].ISTIP == true {
		sumV1 = 0.0
		sumV2 = 0.0000000000001 //0.0
		sumV3 = 0.0
		for i := range tree.CHLD[0].CONTRT {
			x1 = tree.CHLD[0].CONTRT[i]
			x2 = tree.CHLD[1].CONTRT[i]
			x3 = tree.CHLD[2].CONTRT[i]
			sumV1 += (x2 - x1) * (x2 - x1)
			sumV3 += (x2 - x3) * (x2 - x3)
		}
	} else if sumV3 < 0.0 || tree.CHLD[2].ANC == true && tree.CHLD[2].ISTIP == true {
		sumV1 = 0.0
		sumV2 = 0.0
		sumV3 = 0.0000000000001 //0.0
		for i := range tree.CHLD[0].CONTRT {
			x1 = tree.CHLD[0].CONTRT[i]
			x2 = tree.CHLD[1].CONTRT[i]
			x3 = tree.CHLD[2].CONTRT[i]
			sumV1 += (x3 - x1) * (x3 - x1)
			sumV2 += (x3 - x2) * (x3 - x2)
		}
	}
	if sumV1 != 0.0 {
		sumV1 = sumV1 / fntraits
		sumV1 = sumV1 - (tree.CHLD[0].LSPRNLEN - tree.CHLD[0].LSLEN)
	}
	if sumV2 != 0.0 {
		sumV2 = sumV2 / fntraits
		sumV2 = sumV2 - (tree.CHLD[1].LSPRNLEN - tree.CHLD[1].LSLEN)
	}
	if sumV3 != 0.0 {
		sumV3 = sumV3 / fntraits
		sumV3 = sumV3 - (tree.CHLD[2].LSPRNLEN - tree.CHLD[2].LSLEN)
	}
	if sumV1 < 0. {
		sumV1 = 0.0
	}
	if sumV2 < 0. {
		sumV2 = 0.0
	}
	if sumV3 < 0. {
		sumV3 = 0.0
	}
	tree.CHLD[0].LSLEN = sumV1
	tree.CHLD[1].LSLEN = sumV2
	tree.CHLD[2].LSLEN = sumV3
}
