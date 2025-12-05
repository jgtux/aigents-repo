
from config import State
from langgraph.graph import StateGraph, START, END

from typing import TypedDict


class Node:
    """
    NODE
    to add a node to a graph it needs to be a calleble, normaly a function,
    which recive a State and return a State

    name -> name that the user and the api will use to identify the node
    meta -> the actual tructure of the node ti can be an agent or even a executer for sql comands
    callable -> simply the meta wraped in a lambda function, its necessary for create the nodes dinamicly and add them ton the graph

    """
    def __init__(self, name: str, meta=None, callable=None):
        self.name = name
        self.meta = meta
        self.callable = callable


class GraphCore:
    def __init__(self):
        self.graph=StateGraph(State)
        self.start_node=None
        self.end_node=None
        self.nodes=[]
        self.edges=[]
        self.compiled_graph=None


    def reset_graph(self):
        self.graph=StateGraph(State)


    """SAVING INFORMATION"""
    def save_node(self,new_node):
        try:
            self.nodes.append(new_node)
            return new_node
        except Exception as e:
            print(e)
            print("problem saving node")
            return e

    def save_start_node(self,aux_new_node):

        new_node=self.get_one_node(aux_new_node)
        try:
            self.start_node=new_node
            return new_node
        except Exception as e:
            print(e)
            print("problem saving start node")
            return None

    def save_end_node(self,aux_new_node):
        new_node = self.get_one_node(aux_new_node)
        try:
            self.end_node=new_node
            return new_node
        except Exception as e:
            print(e)
            print("problem saving end node")
            return None

    def save_edge(self,node1,node2):
        try:
            aux_edge=(node1,node2)
            self.edges.append(aux_edge)
            return True
        except Exception as e:
            print(e)
            print("problem saving edge")
            return False

    """ADDING TO GRAPH"""
    def add_start_node(self):
        if(self.start_node is not None):
            self.graph.add_node(self.start_node.name,self.start_node.callable)
            self.graph.add_edge(START,self.start_node.name)
        else:
            print("no start node encontered")

    def add_end_node(self):
        if(self.end_node is not None):
            if(self.start_node != self.end_node):
                self.graph.add_node(self.end_node.name, self.end_node.callable)
            self.graph.add_edge(self.end_node.name,END)
        else:
            print("no end node encontered")

    def add_nodes_to_graph(self):
        for n in self.nodes:
            if(n!=self.start_node and n!=self.end_node):
                self.graph.add_node(n.name,n.callable)

    def add_edges_to_graph(self):
        for n1,n2 in self.edges:
            self.graph.add_edge(n1,n2)


    """COMPILING THE GRAPH"""
    def compile_graph(self):
        self.add_start_node()
        self.add_end_node()
        self.add_nodes_to_graph()
        self.add_edges_to_graph()

        self.compiled_graph= self.graph.compile()

    """USING COMPILED GRAPH"""
    def invoke_compiled_graph(self,question:str):

        if self.compiled_graph is None:
            raise RuntimeError("Grafo ainda não foi compilado. Chame compile_graph() antes de invocar.")

        input_state: State = {
            "input": question,
            "output": "",
            "sql_querry": ""
        }
        answer=self.compiled_graph.invoke(input_state)

        return answer


    """GETTING STUFF"""

    def get_nodes(self):
        names=[]
        for i in self.nodes:
            names.append(i.name)
        return names

    def get_one_node(self,name):
        for i in self.nodes:
            if i.name == name:
                return i

        print("node nao encontrado")
        return None


    def get_edges(self):
        return self.edges


