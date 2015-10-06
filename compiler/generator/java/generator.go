package java

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Workiva/frugal/compiler/generator"
	"github.com/Workiva/frugal/compiler/globals"
	"github.com/Workiva/frugal/compiler/parser"
)

const (
	lang                     = "java"
	defaultOutputDir         = "gen-java"
	tab                      = "\t"
	tabtab                   = tab + tab
	tabtabtab                = tab + tab + tab
	tabtabtabtab             = tab + tab + tab + tab
	tabtabtabtabtab          = tab + tab + tab + tab + tab
	tabtabtabtabtabtab       = tab + tab + tab + tab + tab + tab
	tabtabtabtabtabtabtab    = tab + tab + tab + tab + tab + tab + tab
	tabtabtabtabtabtabtabtab = tab + tab + tab + tab + tab + tab + tab + tab
)

type Generator struct {
	*generator.BaseGenerator
	time time.Time
}

func NewGenerator() generator.MultipleFileGenerator {
	return &Generator{&generator.BaseGenerator{}, globals.Now}
}

func (g *Generator) GetOutputDir(dir string, p *parser.Program) string {
	if pkg, ok := p.Namespaces[lang]; ok {
		path := generator.GetPackageComponents(pkg)
		dir = filepath.Join(append([]string{dir}, path...)...)
	}
	return dir
}

func (g *Generator) DefaultOutputDir() string {
	return defaultOutputDir
}

func (g *Generator) GenerateDependencies(p *parser.Program, dir string) error {
	return nil
}

func (g *Generator) CheckCompile(path string) error {
	// TODO
	return nil
}

func (g *Generator) GenerateFile(name, outputDir string, fileType generator.FileType) (*os.File, error) {
	if fileType == generator.CombinedFile {
		return nil, fmt.Errorf("frugal: Bad file type for Java generator: %s", fileType)
	}
	if fileType == generator.PublishFile {
		return g.CreateFile(strings.Title(name)+"Publisher", outputDir, lang, false)
	}
	return g.CreateFile(strings.Title(name)+"Subscriber", outputDir, lang, false)
}

func (g *Generator) GenerateDocStringComment(file *os.File) error {
	comment := fmt.Sprintf(
		"/**\n"+
			" * Autogenerated by Frugal Compiler (%s)\n"+
			" * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING\n"+
			" */",
		globals.Version)

	_, err := file.WriteString(comment)
	return err
}

func (g *Generator) GeneratePackage(file *os.File, p *parser.Program, scope *parser.Scope) error {
	return nil
}

func (g *Generator) GenerateImports(file *os.File, scope *parser.Scope) error {
	imports := "import com.workiva.frugal.Provider;\n"
	imports += "import com.workiva.frugal.Transport;\n"
	imports += "import com.workiva.frugal.TransportFactory;\n"
	imports += "import com.workiva.frugal.Subscription;\n"
	imports += "import org.apache.thrift.TException;\n"
	imports += "import org.apache.thrift.protocol.*;\n"
	imports += "import org.apache.thrift.TApplicationException;\n\n"
	imports += "import org.apache.thrift.transport.TTransportException;\n\n"
	imports += "import org.apache.thrift.transport.TTransportFactory;\n\n"
	imports += "import javax.annotation.Generated;"
	_, err := file.WriteString(imports)
	return err
}

func (g *Generator) GenerateConstants(file *os.File, name string) error {
	return nil
}

func (g *Generator) GeneratePublisher(file *os.File, scope *parser.Scope) error {
	publisher := ""
	publisher += fmt.Sprintf("@Generated(value = \"Autogenerated by Frugal Compiler (%s)\", "+
		"date = \"%s\")\n", globals.Version, g.time.Format("2006-1-2"))
	publisher += fmt.Sprintf("public class %sPublisher {\n\n", scope.Name)

	publisher += fmt.Sprintf(tab+"private static final String delimiter = \"%s\";\n\n", globals.TopicDelimiter)

	publisher += tab + "private Transport transport;\n"
	publisher += tab + "private TProtocol protocol;\n"
	publisher += tab + "private int seqId;\n\n"

	publisher += fmt.Sprintf(tab+"public %sPublisher(Provider provider) {\n", scope.Name)
	publisher += tabtab + "Provider.Client client = provider.build();\n"
	publisher += tabtab + "transport = client.getTransport();\n"
	publisher += tabtab + "protocol = client.getProtocol();\n"
	publisher += tab + "}\n\n"

	args := ""
	if len(scope.Prefix.Variables) > 0 {
		for _, variable := range scope.Prefix.Variables {
			args = fmt.Sprintf("%sString %s, ", args, variable)
		}
	}
	prefix := ""
	for _, op := range scope.Operations {
		publisher += prefix
		prefix = "\n\n"
		publisher += fmt.Sprintf(tab+"public void publish%s(%s%s req) throws TException {\n", op.Name, args, op.Param)
		publisher += fmt.Sprintf(tabtab+"String op = \"%s\";\n", op.Name)
		publisher += fmt.Sprintf(tabtab+"String prefix = %s;\n", generatePrefixStringTemplate(scope))
		publisher += tabtab + "String topic = String.format(\"%s" + scope.Name + "%s%s\", prefix, delimiter, op);\n"
		publisher += tabtab + "transport.preparePublish(topic);\n"
		publisher += tabtab + "seqId++;\n"
		publisher += tabtab + "protocol.writeMessageBegin(new TMessage(op, TMessageType.CALL, seqId));\n"
		publisher += tabtab + "req.write(protocol);\n"
		publisher += tabtab + "protocol.writeMessageEnd();\n"
		publisher += tabtab + "transport.flush();\n"
		publisher += tab + "}\n"
	}

	publisher += "}"

	_, err := file.WriteString(publisher)
	return err
}

func generatePrefixStringTemplate(scope *parser.Scope) string {
	if len(scope.Prefix.Variables) == 0 {
		if scope.Prefix.String == "" {
			return `""`
		}
		return fmt.Sprintf(`"%s%s"`, scope.Prefix.String, globals.TopicDelimiter)
	}
	template := "String.format(\""
	template += scope.Prefix.Template()
	template += globals.TopicDelimiter + "\", "
	prefix := ""
	for _, variable := range scope.Prefix.Variables {
		template += prefix + variable
		prefix = ", "
	}
	template += ")"
	return template
}

func (g *Generator) GenerateSubscriber(file *os.File, scope *parser.Scope) error {
	subscriber := ""
	subscriber += fmt.Sprintf("@Generated(value = \"Autogenerated by Frugal Compiler (%s)\", "+
		"date = \"%s\")\n", globals.Version, g.time.Format("2006-1-2"))
	subscriber += fmt.Sprintf("public class %sSubscriber {\n\n", scope.Name)

	subscriber += fmt.Sprintf(tab+"private static final String delimiter = \"%s\";\n\n", globals.TopicDelimiter)

	subscriber += tab + "private final Provider provider;\n\n"

	subscriber += fmt.Sprintf(tab+"public %sSubscriber(Provider provider) {\n",
		scope.Name)
	subscriber += tabtab + "this.provider = provider;\n"
	subscriber += tab + "}\n\n"

	args := ""
	if len(scope.Prefix.Variables) > 0 {
		for _, variable := range scope.Prefix.Variables {
			args = fmt.Sprintf("%sString %s, ", args, variable)
		}
	}
	prefix := ""
	for _, op := range scope.Operations {
		subscriber += fmt.Sprintf(tab+"public interface %sHandler {\n", op.Name)
		subscriber += fmt.Sprintf(tabtab+"void on%s(%s req);\n", op.Name, op.Param)
		subscriber += tab + "}\n\n"

		subscriber += prefix
		prefix = "\n\n"
		subscriber += fmt.Sprintf(tab+"public Subscription subscribe%s(%sfinal %sHandler handler) throws TException {\n",
			op.Name, args, op.Name)
		subscriber += fmt.Sprintf(tabtab+"final String op = \"%s\";\n", op.Name)
		subscriber += fmt.Sprintf(tabtab+"String prefix = %s;\n", generatePrefixStringTemplate(scope))
		subscriber += tabtab + "String topic = String.format(\"%s" + scope.Name + "%s%s\", prefix, delimiter, op);\n"
		subscriber += tabtab + "final Provider.Client client = provider.build();\n"
		subscriber += tabtab + "Transport transport = client.getTransport();\n"
		subscriber += tabtab + "transport.subscribe(topic);\n\n"

		subscriber += tabtab + "final Subscription sub = new Subscription(topic, transport);\n"
		subscriber += tabtab + "new Thread(new Runnable() {\n"
		subscriber += tabtabtab + "public void run() {\n"
		subscriber += tabtabtabtab + "while (true) {\n"
		subscriber += tabtabtabtabtab + "try {\n"
		subscriber += tabtabtabtabtabtab + fmt.Sprintf("%s received = recv%s(op, client.getProtocol());\n",
			op.Param, op.Name)
		subscriber += tabtabtabtabtabtab + fmt.Sprintf("handler.on%s(received);\n", op.Name)
		subscriber += tabtabtabtabtab + "} catch (TException e) {\n"
		subscriber += tabtabtabtabtabtab + "if (e instanceof TTransportException) {\n"
		subscriber += tabtabtabtabtabtabtab + "TTransportException transportException = (TTransportException) e;\n"
		subscriber += tabtabtabtabtabtabtab + "if (transportException.getType() == TTransportException.END_OF_FILE) {\n"
		subscriber += tabtabtabtabtabtabtabtab + "return;\n"
		subscriber += tabtabtabtabtabtabtab + "}\n"
		subscriber += tabtabtabtabtabtab + "}\n"
		subscriber += tabtabtabtabtabtab + "e.printStackTrace();\n"
		subscriber += tabtabtabtabtabtab + "sub.signal(e);\n"
		subscriber += tabtabtabtabtabtab + "try {\n"
		subscriber += tabtabtabtabtabtabtab + "sub.unsubscribe();\n"
		subscriber += tabtabtabtabtabtab + "} catch (TTransportException e1) {\n"
		subscriber += tabtabtabtabtabtabtab + "e1.printStackTrace();\n"
		subscriber += tabtabtabtabtabtab + "}\n"
		subscriber += tabtabtabtabtab + "}\n"
		subscriber += tabtabtabtab + "}\n"
		subscriber += tabtabtab + "}\n"
		subscriber += tabtab + "}).start();\n\n"

		subscriber += tabtab + "return sub;\n"
		subscriber += tab + "}\n\n"

		subscriber += tab + fmt.Sprintf("private %s recv%s(String op, TProtocol iprot) throws TException {\n", op.Param, op.Name)
		subscriber += tabtab + "TMessage msg = iprot.readMessageBegin();\n"
		subscriber += tabtab + "if (!msg.name.equals(op)) {\n"
		subscriber += tabtabtab + "TProtocolUtil.skip(iprot, TType.STRUCT);\n"
		subscriber += tabtabtab + "iprot.readMessageEnd();\n"
		subscriber += tabtabtab + "throw new TApplicationException(TApplicationException.UNKNOWN_METHOD);\n"
		subscriber += tabtab + "}\n"
		subscriber += tabtab + fmt.Sprintf("%s req = new %s();\n", op.Param, op.Param)
		subscriber += tabtab + "req.read(iprot);\n"
		subscriber += tabtab + "iprot.readMessageEnd();\n"
		subscriber += tabtab + "return req;\n"
		subscriber += tab + "}"
	}
	subscriber += "\n}"

	_, err := file.WriteString(subscriber)
	return err
}
