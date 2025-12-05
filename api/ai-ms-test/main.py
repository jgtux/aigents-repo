from config import State , Config
from graph_core import GraphCore ,Node
from agents_llm_core import AgentLLMCore
from callable_creator import CallableCreator
from langchain.prompts import ChatPromptTemplate
"""
Programa fodastico para rodar Ia Agents
"""
agent_core=AgentLLMCore()
graph_core=GraphCore()
creator=CallableCreator()

altered_grph=False
ok=True
print("vamos criar agentes :) ")

while (ok):
    print("selecione uma opção")
    print("1-criar agente")
    print("2-ver nodes")
    print("3-conectar nodes")
    print("4-invoke")
    print("5-sair")
    op=int(input("opção: "))

    if(op==1):
        print("para criar o agente precisaremos de algumas informações")
        print("nome - model - prompt")
        nome=input("nome: ")

        print("modelos disponiveis:")
        print(agent_core.get_all_llm_names())

        modelo=input("modelo: ")

        print("digite um prompt para o agente")
        a_prompt=input("prompt: ")

        agent=agent_core.create_agent(llm=modelo,tools=None,prompt_text=a_prompt,name=nome)
        callable_agent=creator.create_callable_agent(agent)

        aux_n = Node(name=nome, meta=agent, callable=callable_agent)
        print("agente criado")
        node=graph_core.save_node(aux_n)
        print("salvo como  node",node)

    elif(op==2):
        print(graph_core.get_nodes())

    elif(op==3):
        altered_grph=True
        graph_core.reset_graph()
        print("nodes - > ",graph_core.get_nodes())
        print("selecione primeiro os inicial e final")

        start_n=input("inicial: ")
        graph_core.save_start_node(start_n)
        #graph_core.add_start_node()

        end_n=input("final: ")
        graph_core.save_end_node(end_n)
        #graph_core.add_end_node()

        if(len(graph_core.get_nodes())==1):
            print("so tem um node ele é o começo e o fim")

        elif (graph_core.get_nodes() is not None):
            while(True):
                print("conecte 2 nodes ")
                print("DIGITE 0 QUANDO QUISER SAIR")
                print("- tenha em mente que a ordem importa a resposta do 1º vai para o 2º -")
                n1=input("primeiro: ")
                if(n1=="0"):
                    break
                n2=input("segundo: ")
                if(n2=="0"):
                    break

                graph_core.save_edge(n1,n2)

                stop=input("deseja para?")
                if(stop=="0"):
                    break


    elif (op == 4):
        print("terminado configurações...")
        if (altered_grph==True or graph_core.compiled_graph is None):
            graph_core.compile_graph()
            altered_grph=False
        print("digite a pergunta/entrada")
        pergunta=input("input: ")

        #graph_core.invoke_compiled_graph(pergunta)

        answer = graph_core.invoke_compiled_graph(pergunta)
        print("estado final do grafo:", answer)

        if isinstance(answer, dict) and "output" in answer:
            print("output:", answer["output"])


    elif(op==5):
        print("saindo...")
        ok=False

    else:
        print("opcao invalida")